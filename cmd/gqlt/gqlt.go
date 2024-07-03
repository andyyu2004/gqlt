package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alecthomas/kong"
	"github.com/golang-jwt/jwt/v5"
	"github.com/movio/bramble"
	"github.com/movio/gqlt"
	"github.com/movio/gqlt/ide"
	"github.com/movio/gqlt/lsp"
)

var args struct {
	Run   runCmd   `cmd:"run" help:"Run a query against a GraphQL server."`
	Serve serveCmd `cmd:"serve" default:"1" help:"Run the gqlt lsp server using stdio"`
}

type runCmd struct {
	URL     string            `arg:"" help:"URL of the GraphQL server to query"`
	File    string            `arg:"" help:"Path to the file to run (use - for stdin)"`
	Headers map[string]string `short:"H" help:"HTTP headers to send with the request"`
}

type serveCmd struct{}

func (s *serveCmd) Run() error {
	ide := ide.New()
	server := lsp.New(ide)
	return server.RunStdio()
}

func main() {
	ctx := kong.Parse(&args)
	if err := ctx.Run(); err != nil {
		fmt.Println(err.Error())
	}
}

func (r *runCmd) Run() error {
	file := r.File
	if r.File == "-" {
		file = "/dev/stdin"
	}

	executor := gqlt.New()
	ctx := context.Background()

	headers := make(http.Header)
	for k, v := range r.Headers {
		headers.Add(k, v)
	}

	refreshToken, ok := os.LookupEnv("BRAMBLE_TOKEN")
	if !ok {
		return fmt.Errorf("BRAMBLE_TOKEN environment variable not set")
	}

	rt := gqlt.HTTPRoundTripper{
		RoundTripper: &AuthenticatedClient{
			refreshToken: refreshToken,
			url:          r.URL,
			client:       http.DefaultClient,
			gqlClient:    bramble.NewClient(),
		},
		URL:     r.URL,
		Headers: headers,
	}

	return executor.RunFile(ctx, rt, file)
}

type AuthenticatedClient struct {
	lock         sync.RWMutex
	refreshToken string
	expiresAt    int64
	accessToken  string
	url          string
	client       *http.Client
	gqlClient    *bramble.GraphQLClient
}

func (c *AuthenticatedClient) getAccessToken(ctx context.Context) (string, error) {
	if time.Until(time.Unix(atomic.LoadInt64(&c.expiresAt), 0)) > time.Minute {
		c.lock.RLock()
		defer c.lock.RUnlock()
		return c.accessToken, nil
	}

	return c.refreshAccessToken(ctx)
}

func (c *AuthenticatedClient) refreshAccessToken(ctx context.Context) (string, error) {
	// lock early to prevent multiple requests going out
	c.lock.Lock()
	defer c.lock.Unlock()

	// recheck condition in case another thread has already refreshed the token
	if time.Until(time.Unix(c.expiresAt, 0)) > time.Minute {
		return c.accessToken, nil
	}

	var out struct {
		Bramble struct {
			AccessToken string
		}
	}

	err := c.gqlClient.Request(ctx, c.url, &bramble.Request{
		Query: `mutation ($refreshToken: String!) {
			bramble {
				accessToken: requestAccessToken(refreshToken: $refreshToken)
			}
		}`,
		Variables: map[string]interface{}{
			"refreshToken": c.refreshToken,
		},
	}, &out)
	if err != nil {
		return "", fmt.Errorf("error in access token query: %w", err)
	}

	var claims jwt.RegisteredClaims
	_, _, err = new(jwt.Parser).ParseUnverified(out.Bramble.AccessToken, &claims)
	if err != nil {
		return "", fmt.Errorf("error parsing JWT: %w", err)
	}

	c.accessToken = out.Bramble.AccessToken
	atomic.StoreInt64(&c.expiresAt, claims.ExpiresAt.Unix())

	return c.accessToken, nil
}

func (c *AuthenticatedClient) RoundTrip(req *http.Request) (*http.Response, error) {
	accessToken, err := c.getAccessToken(req.Context())
	if err != nil {
		return nil, fmt.Errorf("unable to get access token: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	return c.client.Do(req)
}

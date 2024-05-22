package eval

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/errors"
	"github.com/movio/gqlt/internal/slice"
)

type Client interface {
	Request(ctx context.Context, req Request, out any) (GraphQLErrors, error)
}

type ClientFactory interface {
	CreateClient(t testing.TB) (context.Context, Client)
}

type ClientFactoryFunc func(t testing.TB) (context.Context, Client)

func (f ClientFactoryFunc) CreateClient(t testing.TB) (context.Context, Client) {
	return f(t)
}

type Request struct {
	Query         string         `json:"query"`
	OperationName string         `json:"operationName,omitempty"`
	Variables     map[string]any `json:"variables,omitempty"`
}

type Response struct {
	Errors GraphQLErrors `json:"errors"`
	Data   any           `json:"data"`
}

type GraphQLErrors []*errors.QueryError

func (errs GraphQLErrors) catch() any {
	// keep in sync with typecheck/expr.go (tryExpr)
	return slice.Map(errs, func(err *errors.QueryError) any {
		return map[string]any{
			"message": err.Message,
			"path":    err.Path,
		}
	})
}

var _ catchable = GraphQLErrors{}

func (e GraphQLErrors) Error() string {
	var buf bytes.Buffer
	for _, err := range e {
		buf.WriteString(err.Error())
		buf.WriteString("\n")
	}
	return buf.String()
}

type GraphQLGophersClient struct {
	Schema *graphql.Schema
}

var _ Client = GraphQLGophersClient{}

func (a GraphQLGophersClient) Request(ctx context.Context, req Request, out any) (GraphQLErrors, error) {
	res := a.Schema.Exec(ctx, req.Query, req.OperationName, req.Variables)

	if err := json.Unmarshal([]byte(res.Data), out); err != nil {
		return nil, err
	}

	return res.Errors, nil
}

type HTTPClient struct {
	Client  *http.Client
	URL     string
	Headers http.Header
}

var _ Client = HTTPClient{}

type clientRoundTripper struct {
	Client *http.Client
}

func (c clientRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}

func (c HTTPClient) Request(ctx context.Context, req Request, out any) (GraphQLErrors, error) {
	return HTTPRoundTripper{
		RoundTripper: clientRoundTripper{Client: c.Client},
		URL:          c.URL,
		Headers:      c.Headers,
	}.Request(ctx, req, out)
}

type HTTPRoundTripper struct {
	RoundTripper http.RoundTripper
	Headers      http.Header
	URL          string
}

var _ Client = HTTPRoundTripper{}

func (c HTTPRoundTripper) Request(ctx context.Context, req Request, out any) (GraphQLErrors, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(req)
	if err != nil {
		return nil, fmt.Errorf("unable to encode request body: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, &buf)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}

	if c.Headers != nil {
		httpReq.Header = c.Headers.Clone()
	}

	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	httpReq.Header.Set("Accept", "application/json; charset=utf-8")

	httpRes, err := c.RoundTripper.RoundTrip(httpReq)
	if err != nil {
		return nil, err
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", httpRes.StatusCode)
	}
	defer httpRes.Body.Close()

	res := Response{Data: out}
	if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Errors, nil
}

type HTTPHandler struct {
	Handler http.Handler
	Headers http.Header
	URL     string
}

var _ Client = GraphQLGophersClient{}

type handlerToRoundTripper struct {
	Handler http.Handler
}

func (h handlerToRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	h.Handler.ServeHTTP(w, req)
	return w.Result(), nil
}

func (c HTTPHandler) Request(ctx context.Context, req Request, out any) (GraphQLErrors, error) {
	return HTTPRoundTripper{
		RoundTripper: handlerToRoundTripper{Handler: c.Handler},
		Headers:      c.Headers,
		URL:          c.URL,
	}.Request(ctx, req, out)
}

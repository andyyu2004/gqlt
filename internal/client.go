package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/graph-gophers/graphql-go"
)

type Client interface {
	Request(ctx context.Context, req Request, out any) error
}

type Request struct {
	Query         string         `json:"query"`
	OperationName string         `json:"operationName,omitempty"`
	Variables     map[string]any `json:"variables,omitempty"`
}

type GraphQLGophersClient struct {
	Schema *graphql.Schema
}

var _ Client = GraphQLGophersClient{}

func (a GraphQLGophersClient) Request(ctx context.Context, req Request, out any) error {
	res := a.Schema.Exec(ctx, req.Query, req.OperationName, req.Variables)
	if len(res.Errors) > 0 {
		errs := make([]error, 0, len(res.Errors))
		for _, err := range res.Errors {
			errs = append(errs, err)
		}

		return errors.Join(errs...)
	}

	return json.Unmarshal([]byte(res.Data), out)
}

type HTTPClient struct {
	Handler http.Handler
	Headers http.Header
	URL     string
}

var _ Client = GraphQLGophersClient{}

func (c HTTPClient) Request(ctx context.Context, req Request, out any) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(req)
	if err != nil {
		return fmt.Errorf("unable to encode request body: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, &buf)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	if c.Headers != nil {
		httpReq.Header = c.Headers.Clone()
	}

	httpReq.Header.Set("Content-Type", "application/json; charset=utf-8")
	httpReq.Header.Set("Accept", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	c.Handler.ServeHTTP(w, httpReq)

	type Response struct {
		Errors GraphqlErrors `json:"errors"`
		Data   interface{}
	}

	graphqlResponse := Response{Data: out}

	if err := json.NewDecoder(w.Body).Decode(&graphqlResponse); err != nil {
		return err
	}

	if len(graphqlResponse.Errors) > 0 {
		return graphqlResponse.Errors
	}

	return nil
}

type Response struct {
	Errors GraphqlErrors `json:"errors"`
	Data   interface{}
}

type GraphqlErrors []GraphqlError

type GraphqlError struct {
	Message string   `json:"message"`
	Path    ast.Path `json:"path,omitempty"`
}

func (e GraphqlErrors) Error() string {
	var errs []string
	for _, err := range e {
		errs = append(errs, err.Message)
	}
	return strings.Join(errs, ",")
}

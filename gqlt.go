package gqlt

import "github.com/movio/gqlt/internal/eval"

var (
	New              = eval.New
	Discover         = eval.Discover
	Ext              = eval.Ext
	WithGlob         = eval.WithGlob
	TypeCheck        = eval.TypeCheck
	WithSchema       = eval.WithSchema
	WithErrorHandler = eval.WithErrorHandler
)

type (
	ClientFactory        = eval.ClientFactory
	ClientFactoryFunc    = eval.ClientFactoryFunc
	Client               = eval.Client
	GraphQLGophersClient = eval.GraphQLGophersClient
	HTTPHandler          = eval.HTTPHandler
	HTTPClient           = eval.HTTPClient
	HTTPRoundTripper     = eval.HTTPRoundTripper
	GraphQLErrors        = eval.GraphQLErrors
	Request              = eval.Request
	Response             = eval.Response
)

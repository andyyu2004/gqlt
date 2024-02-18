package gqlt

import "github.com/movio/gqlt/internal/eval"

var (
	New        = eval.New
	Discover   = eval.Discover
	Ext        = eval.Ext
	WithGlob   = eval.WithGlob
	TypeCheck  = eval.TypeCheck
	WithSchema = eval.WithSchema
)

type (
	ClientFactory        = eval.ClientFactory
	ClientFactoryFunc    = eval.ClientFactoryFunc
	Client               = eval.Client
	GraphQLGophersClient = eval.GraphQLGophersClient
	HTTPClient           = eval.HTTPClient
)

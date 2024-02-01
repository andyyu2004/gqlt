package gqlt

import "github.com/andyyu2004/gqlt/internal"

var (
	New      = internal.New
	Discover = internal.Discover
	Ext      = internal.Ext
	WithGlob = internal.WithGlob
)

type (
	ClientFactory        = internal.ClientFactory
	ClientFactoryFunc    = internal.ClientFactoryFunc
	Client               = internal.Client
	GraphQLGophersClient = internal.GraphQLGophersClient
	HTTPClient           = internal.HTTPClient
)

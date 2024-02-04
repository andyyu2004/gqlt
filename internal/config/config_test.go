package config_test

import (
	"encoding/json"
	"testing"

	"github.com/andyyu2004/gqlt/internal/config"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestParseConfig(t *testing.T) {
	check := func(content string, unmarshal func([]byte, any) error, expected config.Config) {
		t.Helper()
		var c config.Config
		require.NoError(t, unmarshal([]byte(content), &c))
		require.Equal(t, expected, c)
	}

	expected := config.Config{ProjectConfig: config.ProjectConfig{Schema: "schema.graphql"}}
	check(`schema: schema.graphql`, yaml.Unmarshal, expected)
	check(`{ "schema": "schema.graphql" }`, json.Unmarshal, expected)

	// expected = config.Config{
	// 	Projects: map[string]config.ProjectConfig{
	// 		"project1": {Schema: "schema1.graphql"},
	// 		"project2": {Schema: "schema2.graphql"},
	// 	},
	// }
	//
	// check(`projects:
	//  project1:
	//    schema: schema1.graphql
	//  project2:
	//    schema: schema2.graphql`, yaml.Unmarshal, expected)
	// check(`{ "projects": { "project1": { "schema": "schema1.graphql" }, "project2": { "schema": "schema2.graphql" } } }`, json.Unmarshal, expected)
}

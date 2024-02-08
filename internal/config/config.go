package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/andyyu2004/gqlt/gqlparser"
	"github.com/andyyu2004/gqlt/gqlparser/ast"
	"github.com/andyyu2004/gqlt/syn"
	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v2"
)

// Implements subset of `.graphqlrc` config parsing
// https://the-guild.dev/graphql/config/docs/user/usage
// ```typescript
// type ProjectConfig = {
//  schema: string;
// }
//
// type Config = ProjectConfig | { projects: Record<string, ProjectConfig> };
// ```
// Only the `schema` field is supported currently

type ProjectConfig struct {
	Schema string `yaml:"schema"`
}

type Config struct {
	ProjectConfig `yaml:",inline" json:",inline"`
	// Unsupported for now, we would need a way to select which project to use
	// Maybe `set schema <project>`?
	// Projects      map[string]ProjectConfig `yaml:"projects,omitempty" json:"projects,omitempty"`
}

func LoadSchema(workspaces ...string) (*syn.Schema, error) {
	config, root, err := discover(workspaces)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, nil
	}

	return buildSchema(config, root)
}

func discover(workspaces []string) (*Config, string, error) {
	// Load an appropriate `.graphqlrc` config file.
	// We don't really support multiple workspaces, but we'll just use the first one that works.
	for _, root := range workspaces {
		for _, candidate := range []string{".graphqlrc", ".graphqlrc.yaml", ".graphqlrc.yml", ".graphqlrc.json"} {
			configPath := filepath.Join(root, candidate)
			_, err := os.Stat(configPath)
			if err != nil {
				continue
			}

			f, err := os.ReadFile(configPath)
			if err != nil {
				return nil, "", err
			}

			var config Config
			switch filepath.Ext(configPath) {
			case ".json":
				if err = json.Unmarshal(f, &config); err != nil {
					return nil, "", err
				}
			case ".yaml", ".yml", ".graphqlrc":
				if err = yaml.Unmarshal(f, &config); err != nil {
					return nil, "", err
				}
			default:
				panic("missing file case")
			}

			return &config, root, nil
		}
	}

	return nil, "", nil
}

func buildSchema(config *Config, path string) (*syn.Schema, error) {
	schemaPaths, err := doublestar.FilepathGlob(filepath.Join(path, config.Schema), doublestar.WithFilesOnly())
	if err != nil {
		return nil, err
	}

	if schemaPaths == nil {
		return nil, nil
	}

	sources := []*ast.Source{}
	for _, schemaPath := range schemaPaths {
		contents, err := os.ReadFile(schemaPath)
		if err != nil {
			return nil, err
		}

		sources = append(sources, &ast.Source{
			Name:  schemaPath,
			Input: string(contents),
		})
	}

	return gqlparser.LoadSchema(sources...)
}

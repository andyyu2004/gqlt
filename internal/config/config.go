package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/movio/gqlt/gqlparser"
	"github.com/movio/gqlt/gqlparser/ast"
	"github.com/movio/gqlt/syn"
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
	Projects      map[string]ProjectConfig `yaml:"projects,omitempty" json:"projects,omitempty"`
}

type schemaEntry struct {
	glob   string
	schema *syn.Schema
}

type Schemas struct {
	schemas []schemaEntry
}

func (s *Schemas) ForPath(path string) *syn.Schema {
	for _, entry := range s.schemas {
		if matches, err := doublestar.Match(entry.glob, path); err == nil && matches {
			return entry.schema
		}
	}

	return nil
}

func LoadSchemas(workspaces ...string) (*Schemas, error) {
	config, root, err := discover(workspaces)
	if err != nil {
		return nil, err
	}

	return buildSchemas(config, root)
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

func buildSchemas(config *Config, root string) (*Schemas, error) {
	// Can't use the glob directly as it probably ends with *.graphql
	if config == nil {
		// no .graphqlrc file found, use a reasonable default where we scan for all .graphql files
		schema, err := buildSchema(ProjectConfig{Schema: "**/*.graphql"}, root)
		if err != nil {
			return nil, err
		}

		return &Schemas{schemas: []schemaEntry{{"**/*.gqlt", schema}}}, nil
	}

	if config.Projects == nil {
		schema, err := buildSchema(config.ProjectConfig, root)
		if err != nil {
			return nil, err
		}

		return &Schemas{schemas: []schemaEntry{{"**/*.gqlt", schema}}}, nil
	}

	schemas := make([]schemaEntry, 0, len(config.Projects))
	for _, projectConfig := range config.Projects {
		schema, err := buildSchema(projectConfig, root)
		if err != nil {
			return nil, err
		}

		// change the graphql glob to a gqlt glob
		schemas = append(schemas, schemaEntry{strings.ReplaceAll(filepath.Join(root, projectConfig.Schema), ".graphql", ".gqlt"), schema})
	}

	return &Schemas{schemas: schemas}, nil
}

func buildSchema(config ProjectConfig, root string) (*syn.Schema, error) {
	schemaPaths, err := doublestar.FilepathGlob(filepath.Join(root, config.Schema), doublestar.WithFilesOnly())
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

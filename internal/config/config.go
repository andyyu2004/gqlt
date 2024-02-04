package config

// Implements subset of `.graphqlrc` config parsing
// https://the-guild.dev/graphql/config/docs/user/usage
// ```typescript
// type ProjectConfig = {
//  schema: string;
// }
//
// type Config = ProjectConfig | { projects: Record<string, ProjectConfig> };
// ```

type ProjectConfig struct {
	Schema string `yaml:"schema"`
}

type Config struct {
	Projects map[string]ProjectConfig `yaml:"projects,omitempty" json:"projects,omitempty"`
	Schema   string                   `yaml:"schema,omitempty" json:"schema,omitempty"`
}

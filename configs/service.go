package configs

// Service represents metadata that identifies this microservice instance.
// These values are primarily used to configure the CLI commands, logging
// context, and discovery/telemetry endpoints. When cloning the template for
// a new service, update the corresponding values in configs/config.yaml or
// via environment overrides.
type Service struct {
	Name        string   `mapstructure:"name"`
	Slug        string   `mapstructure:"slug"`
	Description string   `mapstructure:"description"`
	Version     string   `mapstructure:"version"`
	Tags        []string `mapstructure:"tags"`
}

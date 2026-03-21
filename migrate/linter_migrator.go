package migrate

// VintRuleConfig represents a single rule's configuration for vint.yaml output.
type VintRuleConfig struct {
	Severity string         `yaml:"severity,omitempty"`
	Disabled bool           `yaml:"disabled,omitempty"`
	Options  map[string]any `yaml:",inline"`
}

// LinterMigrator translates golangci-lint linter settings to vint rule configs.
type LinterMigrator interface {
	// Name returns the golangci-lint linter name (e.g. "errcheck").
	Name() string

	// MigrateConfig translates golangci-lint linter settings to vint.yaml rule configs.
	// settings is the raw map from linters-settings.<linter> in .golangci.yml (may be nil).
	// Returns map of fullVintRulePath → VintRuleConfig.
	MigrateConfig(settings map[string]any) (map[string]VintRuleConfig, error)
}

package gocyclo

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the gocyclo golangci-lint linter.
// gocyclo computes and checks the cyclomatic complexity of functions.
type Migrator struct{}

func (*Migrator) Name() string {
	return "gocyclo"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// noHighCyclomaticComplexity is always enabled when gocyclo is enabled.
	opts := map[string]any{}

	if settings != nil {
		if v, ok := settings["min-complexity"]; ok {
			// Convert to int64 since the vint rule's Configure expects int64,
			// but YAML v3 decodes integers as int.
			switch n := v.(type) {
			case int:
				opts["min-complexity"] = int64(n)
			case int64:
				opts["min-complexity"] = n
			default:
				opts["min-complexity"] = v
			}
		}
	}

	configs["lint/complexity/noHighCyclomaticComplexity"] = migrate.VintRuleConfig{
		Options: opts,
	}

	return configs, nil
}

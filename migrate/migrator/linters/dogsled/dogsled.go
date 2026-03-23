package dogsled

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the dogsled golangci-lint linter.
// dogsled checks for assignments with too many blank identifiers.
type Migrator struct{}

func (*Migrator) Name() string {
	return "dogsled"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// noExcessiveBlankIdentifiers is always enabled when dogsled is enabled.
	opts := map[string]any{}

	if settings != nil {
		if v, ok := settings["max-blank-identifiers"]; ok {
			opts["max-blank-identifiers"] = v
		}
	}

	configs["lint/style/noExcessiveBlankIdentifiers"] = migrate.VintRuleConfig{
		Options: opts,
	}

	return configs, nil
}

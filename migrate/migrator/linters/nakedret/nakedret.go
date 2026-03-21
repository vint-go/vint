package nakedret

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the nakedret golangci-lint linter.
// nakedret checks that functions with naked returns are not longer than
// a configurable maximum size.
type Migrator struct{}

func (*Migrator) Name() string {
	return "nakedret"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// noNakedReturn is always enabled when nakedret is enabled.
	opts := map[string]any{}

	if settings != nil {
		if v, ok := settings["max-func-lines"]; ok {
			opts["max-func-lines"] = v
		}
	}

	configs["lint/style/noNakedReturn"] = migrate.VintRuleConfig{Options: opts}

	return configs, nil
}

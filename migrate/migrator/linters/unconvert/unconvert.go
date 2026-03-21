package unconvert

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the unconvert golangci-lint linter.
// unconvert detects unnecessary type conversions and supports optional
// fast-math and safe settings.
type Migrator struct{}

func (*Migrator) Name() string {
	return "unconvert"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	opts := map[string]any{}

	if settings != nil {
		if v, ok := settings["fast-math"]; ok {
			opts["fast-math"] = v
		}
		if v, ok := settings["safe"]; ok {
			opts["safe"] = v
		}
	}

	config := migrate.VintRuleConfig{}
	if len(opts) > 0 {
		config.Options = opts
	}

	return map[string]migrate.VintRuleConfig{
		"lint/style/noUnnecessaryConversion": config,
	}, nil
}

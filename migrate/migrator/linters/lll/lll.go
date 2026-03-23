package lll

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the lll golangci-lint linter.
// lll checks that lines do not exceed a configured maximum length.
type Migrator struct{}

func (*Migrator) Name() string {
	return "lll"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// noLineTooLong is always enabled when lll is enabled.
	opts := map[string]any{}

	if settings != nil {
		if v, ok := settings["line-length"]; ok {
			opts["line-length"] = v
		}
		if v, ok := settings["tab-width"]; ok {
			opts["tab-width"] = v
		}
	}

	configs["lint/style/noLineTooLong"] = migrate.VintRuleConfig{
		Options: opts,
	}

	return configs, nil
}

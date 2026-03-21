package err113

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the err113 golangci-lint linter.
// err113 has no configuration options — it checks error handling expressions
// (direct error comparison and dynamic error creation).
type Migrator struct{}

func (*Migrator) Name() string {
	return "err113"
}

func (*Migrator) MigrateConfig(_ map[string]any) (map[string]migrate.VintRuleConfig, error) {
	return map[string]migrate.VintRuleConfig{
		"lint/correctness/noDirectErrorComparison": {},
		"lint/correctness/noDynamicErrors":         {},
	}, nil
}

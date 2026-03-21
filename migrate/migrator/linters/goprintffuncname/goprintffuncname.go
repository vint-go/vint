package goprintffuncname

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the goprintffuncname golangci-lint linter.
// goprintffuncname has no configuration options — it simply checks that
// printf-like functions are named with 'f' at the end.
type Migrator struct{}

func (*Migrator) Name() string {
	return "goprintffuncname"
}

func (*Migrator) MigrateConfig(_ map[string]any) (map[string]migrate.VintRuleConfig, error) {
	return map[string]migrate.VintRuleConfig{
		"lint/style/usePrintfSuffix": {},
	}, nil
}

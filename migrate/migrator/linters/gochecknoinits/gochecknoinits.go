package gochecknoinits

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the gochecknoinits golangci-lint linter.
// gochecknoinits has no configuration options — it simply checks that no
// init functions are present in Go code.
type Migrator struct{}

func (*Migrator) Name() string {
	return "gochecknoinits"
}

func (*Migrator) MigrateConfig(_ map[string]any) (map[string]migrate.VintRuleConfig, error) {
	return map[string]migrate.VintRuleConfig{
		"lint/style/noInitFunction": {},
	}, nil
}

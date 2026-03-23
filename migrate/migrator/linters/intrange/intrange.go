package intrange

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the intrange golangci-lint linter.
// intrange detects C-style for loops and redundant len() in range expressions
// that can use Go 1.22+ integer range syntax. It has no configuration options.
type Migrator struct{}

func (*Migrator) Name() string {
	return "intrange"
}

func (*Migrator) MigrateConfig(_ map[string]any) (map[string]migrate.VintRuleConfig, error) {
	return map[string]migrate.VintRuleConfig{
		"lint/style/noRedundantRangeLen": {},
		"lint/style/useIntegerRange":     {},
	}, nil
}

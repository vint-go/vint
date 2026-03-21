package gocheckcompilerdirectives

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the gocheckcompilerdirectives golangci-lint linter.
// gocheckcompilerdirectives has no configuration options — it checks that go
// compiler directive comments (//go:) are valid. Both original rules
// (noSpaceInDirective, noUnrecognizedDirective) are subsumed by noMalformedDirective.
type Migrator struct{}

func (*Migrator) Name() string {
	return "gocheckcompilerdirectives"
}

func (*Migrator) MigrateConfig(_ map[string]any) (map[string]migrate.VintRuleConfig, error) {
	return map[string]migrate.VintRuleConfig{
		"lint/correctness/noMalformedDirective": {},
	}, nil
}

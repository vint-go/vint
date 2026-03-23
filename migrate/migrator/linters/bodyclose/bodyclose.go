package bodyclose

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the bodyclose golangci-lint linter.
// bodyclose has no configuration options — it simply checks that HTTP
// response bodies are closed.
type Migrator struct{}

func (*Migrator) Name() string {
	return "bodyclose"
}

func (*Migrator) MigrateConfig(_ map[string]any) (map[string]migrate.VintRuleConfig, error) {
	return map[string]migrate.VintRuleConfig{
		"lint/correctness/noUnclosedBodies": {},
	}, nil
}

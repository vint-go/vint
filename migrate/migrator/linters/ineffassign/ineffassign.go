package ineffassign

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the ineffassign golangci-lint linter.
// ineffassign detects ineffectual assignments in Go code.
type Migrator struct{}

func (*Migrator) Name() string {
	return "ineffassign"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// noIneffectualAssignment is always enabled when ineffassign is enabled.
	opts := map[string]any{}

	if settings != nil {
		if v, ok := settings["check-escaping-errors"]; ok {
			opts["check-escaping-errors"] = v
		}
	}

	configs["lint/correctness/noIneffectualAssignment"] = migrate.VintRuleConfig{
		Options: opts,
	}

	return configs, nil
}

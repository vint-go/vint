package unparam

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the unparam golangci-lint linter.
// unparam reports unused function parameters, parameters that always receive
// the same value, and result parameters that always return the same value.
// It supports a check-exported setting that controls whether exported functions
// are analyzed.
type Migrator struct{}

func (*Migrator) Name() string {
	return "unparam"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// All three rules are always enabled when unparam is enabled.
	// They all share the same check-exported option from golangci-lint.
	rules := []string{
		"lint/suspicious/noConstantParameter",
		"lint/suspicious/noConstantResult",
		"lint/suspicious/noUnusedParameter",
	}

	for _, rule := range rules {
		opts := map[string]any{}
		if settings != nil {
			if v, ok := settings["check-exported"]; ok {
				opts["check-exported"] = v
			}
		}
		configs[rule] = migrate.VintRuleConfig{Options: opts}
	}

	return configs, nil
}

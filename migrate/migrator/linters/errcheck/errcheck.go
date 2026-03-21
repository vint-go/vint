package errcheck

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the errcheck golangci-lint linter.
// errcheck has three rules in vint and configuration options for exclusions.
type Migrator struct{}

func (*Migrator) Name() string {
	return "errcheck"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// noUncheckedError is always enabled when errcheck is enabled.
	uncheckedErrorOpts := map[string]any{}

	if settings != nil {
		// Map disable-default-exclusions.
		if v, ok := settings["disable-default-exclusions"]; ok {
			uncheckedErrorOpts["disable-default-exclusions"] = v
		}

		// Map exclude-functions.
		if v, ok := settings["exclude-functions"]; ok {
			uncheckedErrorOpts["exclude-functions"] = v
		}
	}

	configs["lint/correctness/noUncheckedError"] = migrate.VintRuleConfig{
		Options: uncheckedErrorOpts,
	}

	// check-type-assertions → enable noUncheckedTypeAssertion.
	enableTypeAssert := false
	if settings != nil {
		if v, ok := settings["check-type-assertions"]; ok {
			if b, ok := v.(bool); ok {
				enableTypeAssert = b
			}
		}
	}
	if enableTypeAssert {
		configs["lint/correctness/noUncheckedTypeAssertion"] = migrate.VintRuleConfig{}
	}

	// check-blank → enable noBlankErrorAssignment.
	enableBlank := false
	if settings != nil {
		if v, ok := settings["check-blank"]; ok {
			if b, ok := v.(bool); ok {
				enableBlank = b
			}
		}
	}
	if enableBlank {
		configs["lint/correctness/noBlankErrorAssignment"] = migrate.VintRuleConfig{}
	}

	return configs, nil
}

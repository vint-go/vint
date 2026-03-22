package unused

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the unused golangci-lint linter.
// unused detects unused constants, fields, functions, types, and variables.
type Migrator struct{}

func (*Migrator) Name() string {
	return "unused"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// noUnusedConstant accepts generated-is-used.
	constOpts := map[string]any{}
	if settings != nil {
		if v, ok := settings["generated-is-used"]; ok {
			constOpts["generated-is-used"] = v
		}
	}
	configs["lint/correctness/noUnusedConstant"] = migrate.VintRuleConfig{Options: constOpts}

	// noUnusedField accepts field-writes-are-uses, exported-fields-are-used, generated-is-used.
	fieldOpts := map[string]any{}
	if settings != nil {
		if v, ok := settings["field-writes-are-uses"]; ok {
			fieldOpts["field-writes-are-uses"] = v
		}
		if v, ok := settings["exported-fields-are-used"]; ok {
			fieldOpts["exported-fields-are-used"] = v
		}
		if v, ok := settings["generated-is-used"]; ok {
			fieldOpts["generated-is-used"] = v
		}
	}
	configs["lint/correctness/noUnusedField"] = migrate.VintRuleConfig{Options: fieldOpts}

	// noUnusedFunction has no configuration options.
	configs["lint/correctness/noUnusedFunction"] = migrate.VintRuleConfig{}

	// noUnusedType accepts generated-is-used.
	typeOpts := map[string]any{}
	if settings != nil {
		if v, ok := settings["generated-is-used"]; ok {
			typeOpts["generated-is-used"] = v
		}
	}
	configs["lint/correctness/noUnusedType"] = migrate.VintRuleConfig{Options: typeOpts}

	// noUnusedVariable accepts post-statements-are-reads, local-variables-are-used, generated-is-used.
	varOpts := map[string]any{}
	if settings != nil {
		if v, ok := settings["post-statements-are-reads"]; ok {
			varOpts["post-statements-are-reads"] = v
		}
		if v, ok := settings["local-variables-are-used"]; ok {
			varOpts["local-variables-are-used"] = v
		}
		if v, ok := settings["generated-is-used"]; ok {
			varOpts["generated-is-used"] = v
		}
	}
	configs["lint/correctness/noUnusedVariable"] = migrate.VintRuleConfig{Options: varOpts}

	// parameters-are-used: when set to false, enable noUnusedParameter.
	// Default is true (parameters considered used, not checked).
	if settings != nil {
		if v, ok := settings["parameters-are-used"]; ok {
			if b, ok := v.(bool); ok && !b {
				configs["lint/suspicious/noUnusedParameter"] = migrate.VintRuleConfig{}
			}
		}
	}

	return configs, nil
}

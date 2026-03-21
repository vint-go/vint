package nolintlint

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the nolintlint golangci-lint linter.
// nolintlint reports ill-formed or insufficient nolint directives.
type Migrator struct{}

func (*Migrator) Name() string {
	return "nolintlint"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// noNolintNonMachineReadable is always enabled when nolintlint is enabled.
	// It catches "// nolint" (with a space) which is not machine-readable.
	configs["lint/style/noNolintNonMachineReadable"] = migrate.VintRuleConfig{}

	// noNolintParseError is always enabled when nolintlint is enabled.
	// It catches malformed nolint directives.
	configs["lint/correctness/noNolintParseError"] = migrate.VintRuleConfig{}

	// noUnusedNolint is enabled by default (allow-unused defaults to false).
	// Only disable if allow-unused is explicitly set to true.
	enableUnused := true
	if settings != nil {
		if v, ok := settings["allow-unused"]; ok {
			if b, ok := v.(bool); ok && b {
				enableUnused = false
			}
		}
	}
	if enableUnused {
		configs["lint/suspicious/noUnusedNolint"] = migrate.VintRuleConfig{}
	}

	// noNolintWithoutExplanation is gated by require-explanation.
	if settings != nil {
		if v, ok := settings["require-explanation"]; ok {
			if b, ok := v.(bool); ok && b {
				opts := map[string]any{}
				// allow-no-explanation maps to the exclude option.
				if v, ok := settings["allow-no-explanation"]; ok {
					opts["exclude"] = v
				}
				configs["lint/style/noNolintWithoutExplanation"] = migrate.VintRuleConfig{Options: opts}
			}
		}
	}

	// noNolintWithoutSpecificLinter is gated by require-specific.
	if settings != nil {
		if v, ok := settings["require-specific"]; ok {
			if b, ok := v.(bool); ok && b {
				configs["lint/style/noNolintWithoutSpecificLinter"] = migrate.VintRuleConfig{}
			}
		}
	}

	return configs, nil
}

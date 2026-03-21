package whitespace

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the whitespace golangci-lint linter.
// whitespace checks for unnecessary newlines at the start and end of
// functions, if, for, etc. It also optionally checks multi-line function
// signatures and multi-line if conditions.
type Migrator struct{}

func (*Migrator) Name() string {
	return "whitespace"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// noLeadingBlankLine and noTrailingBlankLine are always enabled
	// when the whitespace linter is enabled (no config needed).
	configs["lint/style/noLeadingBlankLine"] = migrate.VintRuleConfig{}
	configs["lint/style/noTrailingBlankLine"] = migrate.VintRuleConfig{}

	// multi-if -> enable noMultiLineIfBreak (gated by boolean setting, default false).
	if settings != nil {
		if v, ok := settings["multi-if"]; ok {
			if b, ok := v.(bool); ok && b {
				configs["lint/style/noMultiLineIfBreak"] = migrate.VintRuleConfig{
					Options: map[string]any{"multi-if": true},
				}
			}
		}
	}

	// multi-func -> enable noMultiLineFuncBreak (gated by boolean setting, default false).
	if settings != nil {
		if v, ok := settings["multi-func"]; ok {
			if b, ok := v.(bool); ok && b {
				configs["lint/style/noMultiLineFuncBreak"] = migrate.VintRuleConfig{
					Options: map[string]any{"multi-func": true},
				}
			}
		}
	}

	return configs, nil
}

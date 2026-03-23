package sloglint

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the sloglint golangci-lint linter.
// sloglint ensures consistent code style when using log/slog.
type Migrator struct{}

func (*Migrator) Name() string {
	return "sloglint"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// no-mixed-args defaults to true in golangci-lint, so it is always enabled
	// unless explicitly set to false.
	enableNoMixedArgs := true
	if settings != nil {
		if v, ok := settings["no-mixed-args"]; ok {
			if b, ok := v.(bool); ok {
				enableNoMixedArgs = b
			}
		}
	}
	if enableNoMixedArgs {
		configs["lint/sloglint/no-mixed-args"] = migrate.VintRuleConfig{}
	}

	// kv-only: gated by boolean setting.
	if settings != nil {
		if v, ok := settings["kv-only"]; ok {
			if b, ok := v.(bool); ok && b {
				configs["lint/sloglint/kv-only"] = migrate.VintRuleConfig{}
			}
		}
	}

	// attr-only: gated by boolean setting.
	if settings != nil {
		if v, ok := settings["attr-only"]; ok {
			if b, ok := v.(bool); ok && b {
				configs["lint/sloglint/attr-only"] = migrate.VintRuleConfig{}
			}
		}
	}

	// no-global: gated by non-empty string, pass mode.
	if settings != nil {
		if v, ok := settings["no-global"]; ok {
			if s, ok := v.(string); ok && s != "" {
				configs["lint/sloglint/no-global"] = migrate.VintRuleConfig{
					Options: map[string]any{"mode": s},
				}
			}
		}
	}

	// context: gated by non-empty string, pass mode.
	if settings != nil {
		if v, ok := settings["context"]; ok {
			if s, ok := v.(string); ok && s != "" {
				configs["lint/sloglint/context-only"] = migrate.VintRuleConfig{
					Options: map[string]any{"mode": s},
				}
			}
		}
	}

	// static-msg: gated by boolean setting.
	if settings != nil {
		if v, ok := settings["static-msg"]; ok {
			if b, ok := v.(bool); ok && b {
				configs["lint/sloglint/static-msg"] = migrate.VintRuleConfig{}
			}
		}
	}

	// msg-style: gated by non-empty string, pass style.
	if settings != nil {
		if v, ok := settings["msg-style"]; ok {
			if s, ok := v.(string); ok && s != "" {
				configs["lint/sloglint/msg-style"] = migrate.VintRuleConfig{
					Options: map[string]any{"style": s},
				}
			}
		}
	}

	// no-raw-keys: gated by boolean setting.
	if settings != nil {
		if v, ok := settings["no-raw-keys"]; ok {
			if b, ok := v.(bool); ok && b {
				configs["lint/sloglint/no-raw-keys"] = migrate.VintRuleConfig{}
			}
		}
	}

	// key-naming-case: gated by non-empty string, pass case.
	if settings != nil {
		if v, ok := settings["key-naming-case"]; ok {
			if s, ok := v.(string); ok && s != "" {
				configs["lint/sloglint/key-naming-case"] = migrate.VintRuleConfig{
					Options: map[string]any{"case": s},
				}
			}
		}
	}

	// forbidden-keys: gated by non-empty list, pass keys.
	if settings != nil {
		if v, ok := settings["forbidden-keys"]; ok {
			if keys, ok := v.([]any); ok && len(keys) > 0 {
				configs["lint/sloglint/forbidden-keys"] = migrate.VintRuleConfig{
					Options: map[string]any{"keys": keys},
				}
			}
		}
	}

	// args-on-sep-lines: gated by boolean setting.
	if settings != nil {
		if v, ok := settings["args-on-sep-lines"]; ok {
			if b, ok := v.(bool); ok && b {
				configs["lint/sloglint/args-on-sep-lines"] = migrate.VintRuleConfig{}
			}
		}
	}

	return configs, nil
}

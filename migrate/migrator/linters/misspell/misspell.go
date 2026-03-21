package misspell

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the misspell golangci-lint linter.
// misspell detects commonly misspelled English words in Go source files.
type Migrator struct{}

func (*Migrator) Name() string {
	return "misspell"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// noMisspelledWords is always enabled when misspell is enabled.
	opts := map[string]any{}

	if settings != nil {
		if v, ok := settings["locale"]; ok {
			opts["locale"] = v
		}

		if v, ok := settings["mode"]; ok {
			opts["mode"] = v
		}

		if v, ok := settings["extra-words"]; ok {
			opts["extra-words"] = v
		}

		if v, ok := settings["ignore-rules"]; ok {
			opts["ignore-rules"] = v
		}
	}

	configs["lint/correctness/noMisspelledWords"] = migrate.VintRuleConfig{
		Options: opts,
	}

	return configs, nil
}

package copyloopvar

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the copyloopvar golangci-lint linter.
// copyloopvar detects unnecessary copies of loop variables (Go 1.22+).
// It has one optional setting: check-alias.
type Migrator struct{}

func (*Migrator) Name() string {
	return "copyloopvar"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// noUnnecessaryLoopVarCopy is always enabled when copyloopvar is enabled.
	opts := map[string]any{}
	if settings != nil {
		if v, ok := settings["check-alias"]; ok {
			opts["check-alias"] = v
		}
	}

	configs["lint/style/noUnnecessaryLoopVarCopy"] = migrate.VintRuleConfig{Options: opts}

	return configs, nil
}

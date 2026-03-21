package errorlint

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the errorlint golangci-lint linter.
// errorlint checks for issues with Go 1.13+ error wrapping: direct error
// comparisons (should use errors.Is), type assertions on errors (should use
// errors.As), and fmt.Errorf calls that don't use %w for wrapping.
type Migrator struct{}

func (*Migrator) Name() string {
	return "errorlint"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// comparison (default: true) — enable noDirectErrorComparison.
	enableComparison := true
	if settings != nil {
		if v, ok := settings["comparison"]; ok {
			if b, ok := v.(bool); ok {
				enableComparison = b
			}
		}
	}
	if enableComparison {
		configs["lint/correctness/noDirectErrorComparison"] = migrate.VintRuleConfig{}
	}

	return configs, nil
}

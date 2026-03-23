package errorlint

import "github.com/vint-go/vint/migrate"

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
		compOpts := map[string]any{}
		if settings != nil {
			if v, ok := settings["allowed-errors"]; ok {
				compOpts["allowed-errors"] = v
			}
			if v, ok := settings["allowed-errors-wildcard"]; ok {
				compOpts["allowed-errors-wildcard"] = v
			}
		}
		configs["lint/correctness/noDirectErrorComparison"] = migrate.VintRuleConfig{
			Options: compOpts,
		}
	}

	// asserts (default: true) — enable useErrorsAs.
	enableAsserts := true
	if settings != nil {
		if v, ok := settings["asserts"]; ok {
			if b, ok := v.(bool); ok {
				enableAsserts = b
			}
		}
	}
	if enableAsserts {
		configs["lint/correctness/useErrorsAs"] = migrate.VintRuleConfig{}
	}

	// errorf (default: true) — enable noNonWrappingErrorf.
	enableErrorf := true
	if settings != nil {
		if v, ok := settings["errorf"]; ok {
			if b, ok := v.(bool); ok {
				enableErrorf = b
			}
		}
	}
	if enableErrorf {
		errorfOpts := map[string]any{}
		if settings != nil {
			if v, ok := settings["errorf-multi"]; ok {
				errorfOpts["errorf-multi"] = v
			}
		}
		configs["lint/correctness/noNonWrappingErrorf"] = migrate.VintRuleConfig{
			Options: errorfOpts,
		}
	}

	return configs, nil
}

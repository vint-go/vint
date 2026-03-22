package goconst

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the goconst golangci-lint linter.
// goconst detects repeated strings/numbers that could be replaced by constants,
// finds string literals matching existing constants, and detects duplicate constant values.
type Migrator struct{}

func (*Migrator) Name() string {
	return "goconst"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// noRepeatedStrings is always enabled when goconst is enabled (core behavior).
	strOpts := map[string]any{}
	if settings != nil {
		if v, ok := settings["min-len"]; ok {
			strOpts["min-length"] = v
		}
		if v, ok := settings["min-occurrences"]; ok {
			strOpts["min-occurrences"] = v
		}
		if v, ok := settings["ignore-strings"]; ok {
			strOpts["ignore-strings"] = v
		}
		if v, ok := settings["ignore-tests"]; ok {
			strOpts["ignore-tests"] = v
		}
		if v, ok := settings["ignore-calls"]; ok {
			strOpts["ignore-calls"] = v
		}
		if v, ok := settings["eval-const-expressions"]; ok {
			strOpts["eval-const-expressions"] = v
		}
	}
	configs["lint/style/noRepeatedStrings"] = migrate.VintRuleConfig{Options: strOpts}

	// match-constant defaults to true in golangci-lint, so useMatchingConstant
	// is enabled unless explicitly set to false.
	enableMatchConstant := true
	if settings != nil {
		if v, ok := settings["match-constant"]; ok {
			if b, ok := v.(bool); ok {
				enableMatchConstant = b
			}
		}
	}
	if enableMatchConstant {
		configs["lint/suspicious/useMatchingConstant"] = migrate.VintRuleConfig{}
	}

	// numbers → enable noRepeatedNumbers (default false in golangci-lint).
	enableNumbers := false
	if settings != nil {
		if v, ok := settings["numbers"]; ok {
			if b, ok := v.(bool); ok {
				enableNumbers = b
			}
		}
	}
	if enableNumbers {
		numOpts := map[string]any{}
		if settings != nil {
			if v, ok := settings["min-occurrences"]; ok {
				numOpts["min-occurrences"] = v
			}
			if v, ok := settings["min"]; ok {
				numOpts["min"] = v
			}
			if v, ok := settings["max"]; ok {
				numOpts["max"] = v
			}
		}
		configs["lint/style/noRepeatedNumbers"] = migrate.VintRuleConfig{Options: numOpts}
	}

	// find-duplicates → enable noDuplicateConstants (default false in golangci-lint).
	enableFindDuplicates := false
	if settings != nil {
		if v, ok := settings["find-duplicates"]; ok {
			if b, ok := v.(bool); ok {
				enableFindDuplicates = b
			}
		}
	}
	if enableFindDuplicates {
		configs["lint/suspicious/noDuplicateConstants"] = migrate.VintRuleConfig{}
	}

	return configs, nil
}

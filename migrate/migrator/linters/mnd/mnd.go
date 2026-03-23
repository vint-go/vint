package mnd

import (
	"fmt"
	"strings"

	"github.com/vint-go/vint/migrate"
)

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the mnd golangci-lint linter.
// mnd (Magic Number Detector) checks for magic numbers in Go code.
// It has 6 checks (argument, assign, case, condition, operation, return)
// that are controlled by the "checks" setting, plus shared options for
// ignored-numbers, ignored-files, and ignored-functions.
type Migrator struct{}

func (*Migrator) Name() string {
	return "mnd"
}

// checkToRule maps golangci-lint mnd check names to vint rule paths.
var checkToRule = map[string]string{
	"argument":  "lint/style/noMagicNumberInArgument",
	"assign":    "lint/style/noMagicNumberInAssignment",
	"case":      "lint/style/noMagicNumberInCase",
	"condition": "lint/style/noMagicNumberInCondition",
	"operation": "lint/style/noMagicNumberInOperation",
	"return":    "lint/style/noMagicNumberInReturn",
}

// allChecks is the default set of checks when none are specified.
var allChecks = []string{"argument", "assign", "case", "condition", "operation", "return"}

// checksWithIgnoredFunctions lists checks whose vint rules support ignored-functions.
var checksWithIgnoredFunctions = map[string]bool{
	"argument": true,
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// Determine which checks are enabled.
	enabledChecks := allChecks
	if settings != nil {
		if v, ok := settings["checks"]; ok {
			enabledChecks = toStringSlice(v)
		}
	}

	// Extract shared options from settings.
	var ignoredNumbers string
	var ignoredFiles string
	var ignoredFunctions string

	if settings != nil {
		if v, ok := settings["ignored-numbers"]; ok {
			ignoredNumbers = joinAnySlice(v)
		}
		if v, ok := settings["ignored-files"]; ok {
			ignoredFiles = joinAnySlice(v)
		}
		if v, ok := settings["ignored-functions"]; ok {
			ignoredFunctions = joinAnySlice(v)
		}
	}

	for _, check := range enabledChecks {
		rulePath, ok := checkToRule[check]
		if !ok {
			continue
		}

		opts := map[string]any{}

		if ignoredNumbers != "" {
			opts["ignored-numbers"] = ignoredNumbers
		}
		if ignoredFiles != "" {
			opts["ignored-files"] = ignoredFiles
		}
		// Only the argument rule supports ignored-functions.
		if ignoredFunctions != "" && checksWithIgnoredFunctions[check] {
			opts["ignored-functions"] = ignoredFunctions
		}

		configs[rulePath] = migrate.VintRuleConfig{Options: opts}
	}

	return configs, nil
}

// toStringSlice converts a value that may be []any or []string to []string.
func toStringSlice(v any) []string {
	switch val := v.(type) {
	case []any:
		result := make([]string, 0, len(val))
		for _, item := range val {
			result = append(result, fmt.Sprintf("%v", item))
		}
		return result
	case []string:
		return val
	default:
		return nil
	}
}

// joinAnySlice converts a value that may be []any or []string to a
// comma-separated string.
func joinAnySlice(v any) string {
	parts := toStringSlice(v)
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, ",")
}

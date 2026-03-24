package migrate

import "fmt"

// presetEntry represents a set of vint rules that a single linter contributes
// to an exclusion preset. When the preset is active and the linter is enabled,
// the listed vint rules are candidates for removal.
type presetEntry struct {
	linter string   // golangci-lint linter name
	rules  []string // vint rule full paths to exclude
}

// exclusionPresets maps each golangci-lint exclusion preset name to the set
// of (linter, rules) entries it suppresses. The std-error-handling preset is
// handled specially because it modifies rule configuration rather than
// removing rules entirely.
//
// Reference: https://golangci-lint.run/usage/false-positives/#exclusion-presets
var exclusionPresets = map[string][]presetEntry{
	"comments": {
		{
			linter: "staticcheck",
			rules: []string{
				"lint/style/usePackageComment",    // ST1000
				"lint/style/useFuncDocPrefix",     // ST1020
				"lint/style/useTypeDocPrefix",     // ST1021
				"lint/style/useVarConstDocPrefix", // ST1022
			},
		},
		{
			linter: "revive",
			rules: []string{
				"lint/style/useExportedComment", // exported (.+) should have comment
				"lint/style/usePackageComments", // package comment should be of the form / should have a package comment
			},
		},
	},
	"common-false-positives": {
		{
			linter: "gosec",
			rules: []string{
				"lint/security/noUnsafePackage",           // G103
				"lint/security/noVariableCommandExecution", // G204
				"lint/security/noFileInclusionViaVariable", // G304
			},
		},
	},
	"legacy": {
		{
			linter: "govet",
			rules: []string{
				"lint/correctness/noInvalidUnsafePointer",   // unsafeptr: possible misuse of unsafe.Pointer
				"lint/correctness/noUnbufferedSignalChannel", // sigchanyzer: should have signature
			},
		},
		{
			linter: "staticcheck",
			rules: []string{
				"lint/correctness/noIneffectiveBreak", // SA4011
			},
		},
		{
			linter: "gosec",
			rules: []string{
				"lint/correctness/noUncheckedError",              // G104
				"lint/security/noPermissiveDirectoryPermissions", // G301
				"lint/security/noPermissiveFilePermissions",      // G302
				"lint/security/noPermissiveOsCreate",             // G307
			},
		},
	},
	// std-error-handling is handled by applyStdErrorHandlingPreset.
}

// stdErrorHandlingExcludeFunctions are the specific functions added to
// noUncheckedError's exclude-functions when the std-error-handling preset
// is active. These are the functions that can be enumerated from the
// golangci-lint regex pattern; broader patterns like .*Close and .*Flush
// cannot be translated to exact function names.
var stdErrorHandlingExcludeFunctions = []string{
	"os.Remove",
	"os.RemoveAll",
	"os.Setenv",
	"os.Unsetenv",
}

// validPresets is the set of recognized preset names.
var validPresets = map[string]bool{
	"comments":              true,
	"common-false-positives": true,
	"legacy":                true,
	"std-error-handling":    true,
}

// applyExclusionPresets processes the active exclusion presets and removes
// or modifies rules in allConfigs accordingly. It returns any warnings
// generated during processing.
//
// A rule is only removed if all linters that contributed it have their
// contributions excluded by presets. This prevents removing a rule that
// was enabled by linter A when only linter B's entry is in the preset.
func applyExclusionPresets(
	presets []string,
	enabledLintersSet map[string]bool,
	allConfigs map[string]VintRuleConfig,
	ruleSources map[string]map[string]bool,
) []string {
	var warnings []string

	for _, preset := range presets {
		if !validPresets[preset] {
			warnings = append(warnings, fmt.Sprintf("unknown exclusion preset %q — ignored", preset))
			continue
		}

		if preset == "std-error-handling" {
			w := applyStdErrorHandlingPreset(enabledLintersSet, allConfigs)
			warnings = append(warnings, w...)
			continue
		}

		entries, ok := exclusionPresets[preset]
		if !ok {
			continue
		}

		for _, entry := range entries {
			if !enabledLintersSet[entry.linter] {
				continue
			}
			for _, rule := range entry.rules {
				if _, exists := allConfigs[rule]; !exists {
					continue
				}
				// Remove this linter from the rule's contributing sources.
				if sources, ok := ruleSources[rule]; ok {
					delete(sources, entry.linter)
					// Only remove the rule if no contributing linters remain.
					if len(sources) == 0 {
						delete(allConfigs, rule)
					}
				}
			}
		}
	}

	return warnings
}

// applyStdErrorHandlingPreset handles the std-error-handling preset by adding
// specific exclude-functions to the noUncheckedError rule configuration.
// The golangci-lint regex pattern covers broad categories (.*Close, .*Flush,
// .*print) that cannot be translated to exact function names, so a warning
// is emitted for those.
func applyStdErrorHandlingPreset(
	enabledLintersSet map[string]bool,
	allConfigs map[string]VintRuleConfig,
) []string {
	if !enabledLintersSet["errcheck"] {
		return nil
	}

	const rule = "lint/correctness/noUncheckedError"
	cfg, exists := allConfigs[rule]
	if !exists {
		return nil
	}

	if cfg.Options == nil {
		cfg.Options = make(map[string]any)
	}

	// Merge new exclude-functions with any existing ones.
	existing := anyToStringSlice(cfg.Options["exclude-functions"])
	existingSet := make(map[string]bool, len(existing))
	for _, f := range existing {
		existingSet[f] = true
	}
	for _, f := range stdErrorHandlingExcludeFunctions {
		if !existingSet[f] {
			existing = append(existing, f)
		}
	}

	// Convert to []any for YAML marshaling consistency.
	funcs := make([]any, len(existing))
	for i, f := range existing {
		funcs[i] = f
	}
	cfg.Options["exclude-functions"] = funcs
	allConfigs[rule] = cfg

	return []string{
		`std-error-handling preset: added specific exclude-functions to noUncheckedError; ` +
			`generic patterns (.*Close, .*Flush, .*print) cannot be exactly translated — ` +
			`review your vint.yaml and add exclude-functions for any Close/Flush/print-like ` +
			`methods you want to ignore`,
	}
}

// anyToStringSlice converts an any value to a []string.
// Handles []any (from YAML) and []string.
func anyToStringSlice(v any) []string {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case []any:
		result := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	case []string:
		return val
	default:
		return nil
	}
}

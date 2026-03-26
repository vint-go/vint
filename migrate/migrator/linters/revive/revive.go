package revive

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the revive golangci-lint linter.
//
// Revive has a set of default rules and many opt-in rules. In golangci-lint v2,
// when an explicit `rules:` list is provided under revive settings, it
// **replaces** the default rule set entirely — only the listed rules are active.
// When no `rules:` list is provided, only revive's default rules are active.
//
// This migrator faithfully replicates that behavior.
type Migrator struct{}

func (*Migrator) Name() string {
	return "revive"
}

// reviveRule maps a revive rule name (as used in golangci-lint config) to its
// vint rule path.
type reviveRule struct {
	name     string // revive rule name (e.g. "blank-imports")
	vintPath string // full vint rule path
}

// defaultRules are the revive rules enabled by default (matching defaults.toml
// and upstream revive's defaultRules). Order follows upstream.
var defaultRules = []reviveRule{
	{"blank-imports", "lint/style/noBlankImport"},
	{"context-as-argument", "lint/style/useContextAsFirstParam"},
	{"context-keys-type", "lint/correctness/noContextKeysType"},
	{"dot-imports", "lint/style/noDotImport"},
	{"error-naming", "lint/style/useErrorNaming"},
	{"error-return", "lint/style/useErrorLastReturn"},
	{"error-strings", "lint/style/noErrorStrings"},
	{"errorf", "lint/style/useErrorf"},
	{"exported", "lint/style/useExportedComment"},
	{"indent-error-flow", "lint/style/useIndentErrorFlow"},
	{"package-comments", "lint/style/usePackageComments"},
	{"range", "lint/style/noRedundantRangeVal"},
	{"receiver-naming", "lint/style/useReceiverNaming"},
	{"superfluous-else", "lint/style/noSuperfluousElse"},
	{"time-naming", "lint/style/useTimeNaming"},
	{"unexported-return", "lint/style/noUnexportedReturn"},
	{"var-declaration", "lint/style/noVarDeclaration"},
	{"var-naming", "lint/style/useVarNaming"},
}

// allRules maps every known revive extractor_rule name to its vint rule path.
// This must be kept in sync with rules_registry.toml entries where linter = "revive".
var allRules = map[string]string{
	"atomic":                        "lint/correctness/noAtomicMisuse",
	"banned-characters":             "lint/style/noBannedCharacters",
	"blank-imports":                 "lint/style/noBlankImport",
	"bool-literal-in-expr":          "lint/style/noBoolLiteralInExpr",
	"call-to-gc":                    "lint/performance/noCallToGC",
	"confusing-naming":              "lint/style/noConfusingNaming",
	"confusing-results":             "lint/style/noConfusingResults",
	"context-keys-type":             "lint/correctness/noContextKeysType",
	"cyclomatic":                    "lint/complexity/noCyclomaticComplexity",
	"deep-exit":                     "lint/correctness/noDeepExit",
	"defer":                         "lint/correctness/noDeferGotcha",
	"dot-imports":                   "lint/style/noDotImport",
	"duplicated-imports":            "lint/style/noDuplicatedImports",
	"empty-lines":                   "lint/style/noEmptyLines",
	"error-strings":                 "lint/style/noErrorStrings",
	"argument-limit":                "lint/complexity/noExcessiveArguments",
	"max-control-nesting":           "lint/complexity/noExcessiveControlNesting",
	"file-length-limit":             "lint/complexity/noExcessiveFileLength",
	"function-length":               "lint/complexity/noExcessiveFunctionLength",
	"function-result-limit":         "lint/complexity/noExcessiveFunctionResults",
	"max-public-structs":            "lint/style/noExcessivePublicStructs",
	"flag-parameter":                "lint/style/noFlagParameter",
	"forbidden-call-in-wg-go":       "lint/correctness/noForbiddenCallInWgGo",
	"cognitive-complexity":           "lint/complexity/noHighCognitiveComplexity",
	"identical-branches":            "lint/suspicious/noIdenticalBranches",
	"identical-ifelseif-branches":   "lint/suspicious/noIdenticalIfElseIfBranches",
	"identical-ifelseif-conditions": "lint/suspicious/noIdenticalIfElseIfConditions",
	"identical-switch-branches":     "lint/suspicious/noIdenticalSwitchBranches",
	"identical-switch-conditions":   "lint/suspicious/noIdenticalSwitchConditions",
	"imports-blocklist":             "lint/correctness/noImportsBlocklist",
	"inefficient-map-lookup":        "lint/performance/noInefficientMapLookup",
	"line-length-limit":             "lint/style/noLineLengthLimit",
	"modifies-parameter":            "lint/style/noModifiedParameter",
	"modifies-value-receiver":       "lint/suspicious/noModifiedValueReceiver",
	"nested-structs":                "lint/style/noNestedStructs",
	"package-directory-mismatch":    "lint/style/noPackageDirectoryMismatch",
	"range-val-address":             "lint/correctness/noRangeValAddress",
	"range-val-in-closure":          "lint/correctness/noRangeValInClosure",
	"redefines-builtin-id":          "lint/suspicious/noRedefinesBuiltinId",
	"redundant-build-tag":           "lint/style/noRedundantBuildTag",
	"redundant-import-alias":        "lint/style/noRedundantImportAlias",
	"range":                         "lint/style/noRedundantRangeVal",
	"redundant-test-main-exit":      "lint/style/noRedundantTestMainExit",
	"string-of-int":                 "lint/correctness/noStringOfInt",
	"superfluous-else":              "lint/style/noSuperfluousElse",
	"unconditional-recursion":       "lint/correctness/noUnconditionalRecursion",
	"unexported-naming":             "lint/style/noUnexportedNaming",
	"unexported-return":             "lint/style/noUnexportedReturn",
	"unhandled-error":               "lint/correctness/noUnhandledError",
	"unnecessary-format":            "lint/performance/noUnnecessaryFormat",
	"unnecessary-if":                "lint/style/noUnnecessaryIf",
	"unnecessary-stmt":              "lint/style/noUnnecessaryStmt",
	"unreachable-code":              "lint/correctness/noUnreachableCodeAfterExit",
	"unchecked-type-assertion":      "lint/correctness/noUnsafeTypeAssertion",
	"unsecure-url-scheme":           "lint/security/noUnsecureUrlScheme",
	"unused-parameter":              "lint/suspicious/noUnusedParameter",
	"unused-receiver":               "lint/suspicious/noUnusedReceiver",
	"useless-break":                 "lint/suspicious/noUselessBreak",
	"useless-fallthrough":           "lint/suspicious/noUselessFallthrough",
	"var-declaration":               "lint/style/noVarDeclaration",
	"waitgroup-by-value":            "lint/correctness/noWaitgroupByValue",
	"add-constant":                  "lint/style/useAddConstant",
	"use-any":                       "lint/style/useAny",
	"comment-spacings":              "lint/style/useCommentSpacings",
	"comments-density":              "lint/style/useCommentsDensity",
	"context-as-argument":           "lint/style/useContextAsFirstParam",
	"if-return":                     "lint/style/useDirectReturn",
	"early-return":                  "lint/style/useEarlyReturn",
	"epoch-naming":                  "lint/style/useEpochNaming",
	"error-return":                  "lint/style/useErrorLastReturn",
	"error-naming":                  "lint/style/useErrorNaming",
	"errorf":                        "lint/style/useErrorf",
	"use-errors-new":                "lint/style/useErrorsNew",
	"exported":                      "lint/style/useExportedComment",
	"file-header":                   "lint/style/useFileHeader",
	"filename-format":               "lint/style/useFilenameFormat",
	"use-fmt-print":                 "lint/style/useFmtPrint",
	"get-return":                    "lint/style/useGetterReturn",
	"import-alias-naming":           "lint/style/useImportAliasNaming",
	"indent-error-flow":             "lint/style/useIndentErrorFlow",
	"enforce-map-style":             "lint/style/useMapStyle",
	"optimize-operands-order":       "lint/performance/useOptimalOperandsOrder",
	"package-comments":              "lint/style/usePackageComments",
	"package-naming":                "lint/style/usePackageNaming",
	"receiver-naming":               "lint/style/useReceiverNaming",
	"enforce-repeated-arg-type-style": "lint/style/useRepeatedArgTypeStyle",
	"enforce-slice-style":           "lint/style/useSliceStyle",
	"use-slices-sort":               "lint/performance/useSlicesSort",
	"string-format":                 "lint/style/useStringFormat",
	"enforce-switch-style":          "lint/style/useSwitchDefaultStyle",
	"time-date":                     "lint/correctness/useTimeDate",
	"time-equal":                    "lint/correctness/useTimeEqual",
	"time-naming":                   "lint/style/useTimeNaming",
	"var-naming":                    "lint/style/useVarNaming",
	"use-waitgroup-go":              "lint/style/useWaitgroupGo",
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// Parse the rules list from revive settings.
	var configuredRules []ruleEntry
	if settings != nil {
		configuredRules = parseRulesList(settings["rules"])
	}

	if len(configuredRules) == 0 {
		// No explicit rules list → enable only revive defaults.
		for _, r := range defaultRules {
			configs[r.vintPath] = migrate.VintRuleConfig{}
		}
		return configs, nil
	}

	// Explicit rules list provided → it replaces the defaults entirely.
	// Only the listed rules are active (unless disabled: true).
	for _, entry := range configuredRules {
		if entry.disabled {
			continue
		}

		vintPath, ok := allRules[entry.name]
		if !ok {
			continue
		}

		cfg := migrate.VintRuleConfig{}
		if len(entry.arguments) > 0 {
			cfg.Options = mapArguments(entry.name, entry.arguments)
		}
		configs[vintPath] = cfg
	}

	return configs, nil
}

// ruleEntry represents a single rule from the golangci-lint revive config.
type ruleEntry struct {
	name      string
	disabled  bool
	arguments []any
}

// parseRulesList extracts rule entries from the YAML-parsed rules value.
// The expected structure is []any where each element is map[string]any
// with keys: name (string), disabled (bool), arguments ([]any).
func parseRulesList(v any) []ruleEntry {
	if v == nil {
		return nil
	}

	list, ok := v.([]any)
	if !ok {
		return nil
	}

	var entries []ruleEntry
	for _, item := range list {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}

		name, _ := m["name"].(string)
		if name == "" {
			continue
		}

		entry := ruleEntry{name: name}

		if d, ok := m["disabled"].(bool); ok {
			entry.disabled = d
		}

		if args, ok := m["arguments"].([]any); ok {
			entry.arguments = args
		}

		entries = append(entries, entry)
	}

	return entries
}

// mapArguments converts revive rule arguments to vint rule options.
// Most revive rules use positional arguments; we map them to named options
// where the vint rule expects specific keys.
func mapArguments(ruleName string, args []any) map[string]any {
	if len(args) == 0 {
		return nil
	}

	opts := make(map[string]any)

	switch ruleName {
	case "argument-limit":
		if v, ok := toInt(args[0]); ok {
			opts["max"] = v
		}
	case "max-control-nesting":
		if v, ok := toInt(args[0]); ok {
			opts["max"] = v
		}
	case "cognitive-complexity":
		if v, ok := toInt(args[0]); ok {
			opts["max-complexity"] = v
		}
	case "cyclomatic":
		if v, ok := toInt(args[0]); ok {
			opts["min-complexity"] = v
		}
	case "max-public-structs":
		if v, ok := toInt(args[0]); ok {
			opts["max"] = v
		}
	case "function-result-limit":
		if v, ok := toInt(args[0]); ok {
			opts["max"] = v
		}
	case "line-length-limit":
		if v, ok := toInt(args[0]); ok {
			opts["line-length"] = v
		}
	case "function-length":
		if len(args) >= 2 {
			if v, ok := toInt(args[0]); ok {
				opts["statements"] = v
			}
			if v, ok := toInt(args[1]); ok {
				opts["lines"] = v
			}
		} else if len(args) == 1 {
			if v, ok := toInt(args[0]); ok {
				opts["statements"] = v
			}
		}
	case "file-length-limit":
		if v, ok := toInt(args[0]); ok {
			opts["max"] = v
		}
	case "import-alias-naming":
		if v, ok := args[0].(string); ok {
			opts["allow-regex"] = v
		}
	case "context-as-argument":
		if m, ok := args[0].(map[string]any); ok {
			for k, v := range m {
				opts[k] = v
			}
		}
	case "var-naming":
		if len(args) >= 1 {
			opts["allowlist"] = args[0]
		}
		if len(args) >= 2 {
			opts["blocklist"] = args[1]
		}
	case "add-constant":
		if m, ok := args[0].(map[string]any); ok {
			for k, v := range m {
				opts[k] = v
			}
		}
	case "unhandled-error":
		opts["exclude-functions"] = args
	case "comment-spacings":
		opts["markers"] = args
	default:
		// For rules we don't have specific mapping for, pass arguments
		// through as-is under "arguments" key.
		opts["arguments"] = args
	}

	return opts
}

// toInt converts an any value to int, handling float64 (from JSON/YAML).
func toInt(v any) (int, bool) {
	switch val := v.(type) {
	case int:
		return val, true
	case int64:
		return int(val), true
	case float64:
		return int(val), true
	}
	return 0, false
}

package govet

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the govet golangci-lint linter.
// govet wraps Go's "go vet" tool and supports enable/disable/enable-all/disable-all
// for individual analyzers. Each analyzer maps to a specific vint rule.
type Migrator struct{}

func (*Migrator) Name() string {
	return "govet"
}

// analyzerToRule maps golangci-lint govet analyzer names to vint rule paths.
var analyzerToRule = map[string]string{
	// Default analyzers (enabled by default in go vet).
	"appends":            "lint/correctness/noSingleArgAppend",
	"asmdecl":            "lint/correctness/noAsmDeclMismatch",
	"assign":             "lint/correctness/noSelfAssignment",
	"atomic":             "lint/correctness/noAtomicAssignMisuse",
	"bools":              "lint/suspicious/noDuplicateSubExpression",
	"buildtag":           "lint/correctness/noMalformedBuildTag",
	"cgocall":            "lint/correctness/noCgoPointerViolation",
	"composite":          "lint/style/noUnkeyedLiteral",
	"copylock":           "lint/correctness/noCopiedLock",
	"defers":             "lint/correctness/noDeferTimeMisuse",
	"directive":          "lint/correctness/noMalformedDirective",
	"errorsas":           "lint/correctness/noInvalidErrorsAs",
	"framepointer":       "lint/correctness/noFramePointerClobber",
	"hostport":           "lint/correctness/useJoinHostPort",
	"httpresponse":       "lint/correctness/noHttpResponseMisuse",
	"ifaceassert":        "lint/correctness/noImpossibleInterfaceAssert",
	"loopclosure":        "lint/correctness/noLoopClosureCapture",
	"lostcancel":         "lint/correctness/noLostCancel",
	"nilfunc":            "lint/correctness/noNilFuncComparison",
	"printf":             "lint/correctness/noPrintfFormatMismatch",
	"shift":              "lint/correctness/noExcessiveShift",
	"sigchanyzer":        "lint/correctness/noUnbufferedSignalChannel",
	"slog":               "lint/correctness/noSlogKeyValueMismatch",
	"stdmethods":         "lint/correctness/noStdMethodSignatureMismatch",
	"stdversion":         "lint/correctness/noStdlibVersionMismatch",
	"stringintconv":      "lint/correctness/noStringIntConversion",
	"structtag":          "lint/correctness/noMalformedStructTag",
	"testinggoroutine":   "lint/correctness/noTestFatalInGoroutine",
	"tests":              "lint/correctness/noMalformedTestFunction",
	"timeformat":         "lint/correctness/noIncorrectTimeFormat",
	"unmarshal":          "lint/correctness/noNonPointerUnmarshal",
	"unreachable":        "lint/correctness/noUnreachableCode",
	"unsafeptr":          "lint/correctness/noInvalidUnsafePointer",
	"unusedresult":       "lint/correctness/noUnusedFunctionResult",
	"waitgroup":          "lint/correctness/noWaitGroupMisuse",

	// Non-default analyzers (disabled by default, must be explicitly enabled).
	"atomicalign":         "lint/correctness/noAtomicAlignmentIssue",
	"deepequalerrors":     "lint/suspicious/noDeepEqualErrors",
	"fieldalignment":      "lint/performance/useOptimalFieldAlignment",
	"findcall":            "lint/suspicious/noSpecificFunctionCall",
	"httpmux":             "lint/correctness/noConflictingHttpMuxPatterns",
	"nilness":             "lint/correctness/noNilDereference",
	"reflectvaluecompare": "lint/suspicious/noReflectValueCompare",
	"shadow":              "lint/suspicious/noVariableShadowing",
	"sortslice":           "lint/correctness/noInvalidSortSliceArg",
	"unusedwrite":         "lint/correctness/noUnusedWrite",
}

// defaultAnalyzers is the set of analyzers enabled by default in go vet / golangci-lint govet.
var defaultAnalyzers = map[string]bool{
	"appends":          true,
	"asmdecl":          true,
	"assign":           true,
	"atomic":           true,
	"bools":            true,
	"buildtag":         true,
	"cgocall":          true,
	"composite":        true,
	"copylock":         true,
	"defers":           true,
	"directive":        true,
	"errorsas":         true,
	"framepointer":     true,
	"hostport":         true,
	"httpresponse":     true,
	"ifaceassert":      true,
	"loopclosure":      true,
	"lostcancel":       true,
	"nilfunc":          true,
	"printf":           true,
	"shift":            true,
	"sigchanyzer":      true,
	"slog":             true,
	"stdmethods":       true,
	"stdversion":       true,
	"stringintconv":    true,
	"structtag":        true,
	"testinggoroutine": true,
	"tests":            true,
	"timeformat":       true,
	"unmarshal":        true,
	"unreachable":      true,
	"unsafeptr":        true,
	"unusedresult":     true,
	"waitgroup":        true,
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// Determine which analyzers are enabled based on govet settings.
	enableAll := false
	disableAll := false
	var enable []string
	var disable []string

	if settings != nil {
		if v, ok := settings["enable-all"]; ok {
			if b, ok := v.(bool); ok {
				enableAll = b
			}
		}
		if v, ok := settings["disable-all"]; ok {
			if b, ok := v.(bool); ok {
				disableAll = b
			}
		}
		if v, ok := settings["enable"]; ok {
			if list, ok := v.([]any); ok {
				for _, item := range list {
					if s, ok := item.(string); ok {
						enable = append(enable, s)
					}
				}
			}
		}
		if v, ok := settings["disable"]; ok {
			if list, ok := v.([]any); ok {
				for _, item := range list {
					if s, ok := item.(string); ok {
						disable = append(disable, s)
					}
				}
			}
		}
	}

	enableSet := make(map[string]bool, len(enable))
	for _, name := range enable {
		enableSet[name] = true
	}
	disableSet := make(map[string]bool, len(disable))
	for _, name := range disable {
		disableSet[name] = true
	}

	// Replicate golangci-lint's analyzer selection logic.
	for analyzer, vintRule := range analyzerToRule {
		enabled := isAnalyzerEnabled(analyzer, enableAll, disableAll, enableSet, disableSet)
		if !enabled {
			continue
		}

		configs[vintRule] = migrate.VintRuleConfig{}
	}

	// Handle per-analyzer settings if present.
	if settings != nil {
		if perAnalyzer, ok := settings["settings"]; ok {
			if settingsMap, ok := perAnalyzer.(map[string]any); ok {
				// Map findcall settings (name option).
				if findcallSettings, ok := settingsMap["findcall"]; ok {
					if fc, ok := findcallSettings.(map[string]any); ok {
						if _, exists := configs["lint/suspicious/noSpecificFunctionCall"]; exists {
							opts := map[string]any{}
							if name, ok := fc["name"]; ok {
								opts["name"] = name
							}
							if len(opts) > 0 {
								configs["lint/suspicious/noSpecificFunctionCall"] = migrate.VintRuleConfig{
									Options: opts,
								}
							}
						}
					}
				}

				// Map printf settings (funcs option).
				if printfSettings, ok := settingsMap["printf"]; ok {
					if pf, ok := printfSettings.(map[string]any); ok {
						if _, exists := configs["lint/correctness/noPrintfFormatMismatch"]; exists {
							opts := map[string]any{}
							if funcs, ok := pf["funcs"]; ok {
								opts["funcs"] = funcs
							}
							if len(opts) > 0 {
								configs["lint/correctness/noPrintfFormatMismatch"] = migrate.VintRuleConfig{
									Options: opts,
								}
							}
						}
					}
				}

				// Map shadow settings (strict option).
				if shadowSettings, ok := settingsMap["shadow"]; ok {
					if ss, ok := shadowSettings.(map[string]any); ok {
						if _, exists := configs["lint/suspicious/noVariableShadowing"]; exists {
							opts := map[string]any{}
							if strict, ok := ss["strict"]; ok {
								opts["strict"] = strict
							}
							if len(opts) > 0 {
								configs["lint/suspicious/noVariableShadowing"] = migrate.VintRuleConfig{
									Options: opts,
								}
							}
						}
					}
				}
			}
		}
	}

	return configs, nil
}

// isAnalyzerEnabled determines whether a given analyzer should be enabled,
// mirroring golangci-lint's govet selection logic.
func isAnalyzerEnabled(name string, enableAll, disableAll bool, enableSet, disableSet map[string]bool) bool {
	switch {
	case enableAll:
		return !disableSet[name]
	case enableSet[name]:
		return true
	case disableSet[name]:
		return false
	case disableAll:
		return false
	default:
		return defaultAnalyzers[name]
	}
}

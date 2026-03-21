package staticcheck

import (
	"strings"
	"unicode"

	"github.com/strowk/vint/migrate"
)

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the staticcheck golangci-lint linter.
// staticcheck combines multiple check suites: SA (staticcheck), S (simple/gosimple),
// ST (stylecheck), and QF (quickfix). It supports a "checks" setting with
// glob patterns and negation to control which checks are enabled.
type Migrator struct{}

func (*Migrator) Name() string {
	return "staticcheck"
}

// checkToRule maps staticcheck check IDs to their vint rule paths.
// Only checks that have corresponding vint rule implementations are included.
var checkToRule = map[string]string{
	// SA — staticcheck (correctness and suspicious code)
	"SA1000": "lint/correctness/noInvalidRegexp",
	"SA1001": "lint/correctness/noInvalidTemplate",
	"SA1003": "lint/correctness/noInvalidBinaryArg",
	"SA1004": "lint/suspicious/noSuspiciousTimeSleep",
	"SA1005": "lint/correctness/noInvalidExecCommandArg",
	"SA1007": "lint/correctness/noInvalidUrlParse",
	"SA1008": "lint/correctness/noNonCanonicalHeaderKey",
	"SA1011": "lint/correctness/noInvalidUtf8StringArg",
	"SA1012": "lint/correctness/noNilContext",
	"SA1015": "lint/correctness/noTimeTick",
	"SA1016": "lint/correctness/noUntrappableSignal",
	"SA1019": "lint/correctness/noDeprecatedUsage",
	"SA1020": "lint/correctness/noInvalidHostPort",
	"SA1021": "lint/correctness/noNetIpBytesEqual",
	"SA1023": "lint/correctness/noWriterBufferModification",
	"SA1024": "lint/suspicious/noDuplicateCutsetChars",
	"SA1025": "lint/correctness/noTimerResetRetval",
	"SA1026": "lint/correctness/noUnmarshalableType",
	"SA1029": "lint/correctness/noInappropriateContextKey",
	"SA1030": "lint/correctness/noInvalidStrconvArg",
	"SA1031": "lint/correctness/noOverlappingEncoderSlice",
	"SA2001": "lint/suspicious/noEmptyCriticalSection",
	"SA3000": "lint/correctness/noTestMainWithoutExit",
	"SA3001": "lint/correctness/noBenchmarkNAssignment",
	"SA4001": "lint/suspicious/noAddressOfDereference",
	"SA4003": "lint/suspicious/noUnsignedNegativeComparison",
	"SA4004": "lint/suspicious/noSingleIterationLoop",
	"SA4005": "lint/suspicious/noUnobservedFieldAssign",
	"SA4008": "lint/suspicious/noInvariantLoopCondition",
	"SA4009": "lint/suspicious/noOverwrittenArgument",
	"SA4010": "lint/correctness/noDiscardedAppend",
	"SA4011": "lint/correctness/noIneffectiveBreak",
	"SA4014": "lint/suspicious/noDuplicateIfCondition",
	"SA4015": "lint/suspicious/noUselessMathCall",
	"SA4016": "lint/suspicious/noIneffectiveBitwiseOp",
	"SA4019": "lint/correctness/noDuplicateBuildConstraint",
	"SA4022": "lint/suspicious/noAddressNilComparison",
	"SA4023": "lint/suspicious/noImpossibleNilComparison",
	"SA4024": "lint/suspicious/noImpossibleBuiltinResult",
	"SA4025": "lint/suspicious/noIntegerDivisionTruncation",
	"SA4026": "lint/suspicious/noImpreciseConstant",
	"SA4027": "lint/suspicious/noIgnoredQueryModification",
	"SA4028": "lint/suspicious/noModuloOne",
	"SA4030": "lint/suspicious/noIneffectiveRandCall",
	"SA4031": "lint/suspicious/noNeverNilCheck",
	"SA4012": "lint/correctness/noNanComparison",
	"SA5000": "lint/correctness/noNilMapAssignment",
	"SA5001": "lint/correctness/noDeferCloseBeforeErrCheck",
	"SA5002": "lint/correctness/noEmptyForLoop",
	"SA5003": "lint/correctness/noDeferInInfiniteLoop",
	"SA5004": "lint/correctness/noSelectBreakConfusion",
	"SA5005": "lint/correctness/noSelfReferencingFinalizer",
	"SA5007": "lint/correctness/noInfiniteRecursion",
	"SA5012": "lint/correctness/noOddSizeSliceArg",
	"SA6000": "lint/performance/useCompiledRegexp",
	"SA6001": "lint/performance/useStringMapKey",
	"SA6002": "lint/performance/usePointerInSyncPool",
	"SA6003": "lint/performance/noRedundantRuneConversion",
	"SA9002": "lint/suspicious/noNonOctalFileMode",
	"SA9003": "lint/suspicious/noEmptyBranch",
	"SA9004": "lint/suspicious/noImplicitConstValue",
	"SA9005": "lint/suspicious/noUnmarshalableStruct",
	"SA9006": "lint/suspicious/noDubiousBitShift",
	"SA9007": "lint/suspicious/noTempDirDeletion",
	"SA9008": "lint/suspicious/noTypeAssertElseMisread",

	// S — simple/gosimple (simplification)
	"S1000": "lint/style/noSingleCaseSelect",
	"S1001": "lint/style/useCopyBuiltin",
	"S1002": "lint/style/noExplicitBoolComparison",
	"S1003": "lint/style/useStringsContains",
	"S1004": "lint/style/useBytesEqual",
	"S1005": "lint/style/noUnnecessaryBlankIdentifier",
	"S1006": "lint/style/useInfiniteFor",
	"S1007": "lint/style/useRawStringRegexp",
	"S1008": "lint/style/useSimplifiedBoolReturn",
	"S1009": "lint/style/noRedundantNilSliceCheck",
	"S1010": "lint/style/noDefaultSliceIndex",
	"S1011": "lint/style/useSliceAppend",
	"S1012": "lint/style/useTimeSince",
	"S1016": "lint/style/useTypeConversion",
	"S1017": "lint/style/useTrimFunction",
	"S1018": "lint/style/useCopyForSlide",
	"S1019": "lint/style/noRedundantMakeArgs",
	"S1020": "lint/style/noRedundantNilTypeCheck",
	"S1021": "lint/style/useMergedVarDecl",
	"S1023": "lint/style/noRedundantControlFlow",
	"S1024": "lint/style/useTimeUntil",
	"S1028": "lint/style/useErrorMethod",
	"S1029": "lint/style/useDirectStringRange",
	"S1030": "lint/style/useBufferStringOrBytes",
	"S1031": "lint/style/noRedundantNilLoopCheck",
	"S1033": "lint/style/noGuardAroundDelete",
	"S1034": "lint/style/useTypeAssertResult",
	"S1035": "lint/style/noRedundantCanonicalHeaderKey",
	"S1036": "lint/style/noGuardAroundMapAccess",
	"S1037": "lint/style/useTimeSleep",
	"S1038": "lint/style/useSimplifiedPrintFormat",

	// ST — stylecheck
	"ST1000": "lint/style/usePackageComment",
	"ST1001": "lint/style/noDotImport",
	"ST1003": "lint/style/useIdiomaticNaming",
	"ST1005": "lint/style/noCapitalizedErrorString",
	"ST1006": "lint/style/useIdiomaticReceiverName",
	"ST1008": "lint/style/useErrorLastReturn",
	"ST1011": "lint/style/useIdiomaticDurationName",
	"ST1012": "lint/style/useIdiomaticErrorName",
	"ST1016": "lint/style/useConsistentReceiverName",
	"ST1018": "lint/style/noControlCharInString",
	"ST1020": "lint/style/useFuncDocPrefix",
	"ST1021": "lint/style/useTypeDocPrefix",
	"ST1022": "lint/style/useVarConstDocPrefix",
	"ST1023": "lint/style/noRedundantVarType",

	// QF — quickfix
	"QF1002": "lint/style/useTaggedSwitch",
	"QF1005": "lint/performance/useInlineMathPow",
	"QF1006": "lint/complexity/useLoopCondition",
	"QF1007": "lint/style/useMergedConditionalDecl",
	"QF1008": "lint/style/useSimplifiedSelector",
	"QF1009": "lint/correctness/useTimeEqual",
	"QF1010": "lint/style/useStringConversionInPrint",
}

// allCheckIDs lists every check ID known to the mapping.
// This is used by the filter to resolve "all" and glob patterns.
var allCheckIDs []string

func init() {
	for id := range checkToRule {
		allCheckIDs = append(allCheckIDs, id)
	}
	sortStrings(allCheckIDs)
}

// defaultChecks mirrors golangci-lint's default staticcheck checks setting.
var defaultChecks = []string{"all", "-ST1000", "-ST1003", "-ST1016", "-ST1020", "-ST1021", "-ST1022"}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// Determine which checks are requested.
	checks := defaultChecks
	if settings != nil {
		if v, ok := settings["checks"]; ok {
			if list := toStringSlice(v); len(list) > 0 {
				checks = list
			}
		}
	}

	// Filter checks using staticcheck's pattern matching logic.
	enabled := filterChecks(allCheckIDs, checks)

	// Build options for rules that accept settings.
	var dotImportWhitelist []string
	if settings != nil {
		if v, ok := settings["dot-import-whitelist"]; ok {
			dotImportWhitelist = toStringSlice(v)
		}
	}

	for _, checkID := range allCheckIDs {
		if !enabled[checkID] {
			continue
		}
		vintPath, ok := checkToRule[checkID]
		if !ok {
			continue
		}

		cfg := migrate.VintRuleConfig{}

		// Pass dot-import-whitelist to the noDotImport rule as allowedPackages.
		if checkID == "ST1001" && len(dotImportWhitelist) > 0 {
			cfg.Options = map[string]any{
				"allowedPackages": dotImportWhitelist,
			}
		}

		configs[vintPath] = cfg
	}

	return configs, nil
}

// filterChecks replicates staticcheck's check selection logic.
// It processes patterns in order: "all" enables everything,
// "-SA1000" disables a specific check, "SA*" enables a category, etc.
func filterChecks(available []string, checks []string) map[string]bool {
	allowed := make(map[string]bool)

	for _, check := range checks {
		enable := true
		if len(check) > 1 && check[0] == '-' {
			enable = false
			check = check[1:]
		}

		if check == "*" || check == "all" {
			for _, c := range available {
				allowed[c] = enable
			}
		} else if strings.HasSuffix(check, "*") {
			prefix := check[:len(check)-1]
			// Determine if this is a category-level glob (letters only)
			// or a more specific prefix glob (letters + numbers).
			isCat := strings.IndexFunc(prefix, func(r rune) bool {
				return unicode.IsNumber(r)
			}) == -1

			for _, a := range available {
				idx := strings.IndexFunc(a, func(r rune) bool {
					return unicode.IsNumber(r)
				})
				if isCat {
					if idx > 0 {
						cat := a[:idx]
						if prefix == cat {
							allowed[a] = enable
						}
					}
				} else {
					if strings.HasPrefix(a, prefix) {
						allowed[a] = enable
					}
				}
			}
		} else {
			allowed[check] = enable
		}
	}

	return allowed
}

func toStringSlice(v any) []string {
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

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

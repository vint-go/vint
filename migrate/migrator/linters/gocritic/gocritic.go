package gocritic

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the gocritic golangci-lint linter.
// gocritic is a large linter with many checks covering correctness,
// performance, style, and complexity. Only non-experimental,
// non-opinionated checks are enabled by default. The enabled set
// can be customized via enabled-tags, disabled-tags, enabled-checks,
// and disabled-checks settings in the golangci-lint config.
type Migrator struct{}

func (*Migrator) Name() string {
	return "gocritic"
}

// gocriticCheck maps a gocritic checker to its vint rule and tags.
type gocriticCheck struct {
	checker  string   // gocritic checker name
	vintRule string   // full vint rule path
	tags     []string // gocritic tags (diagnostic/style/performance + experimental/opinionated)
}

// allChecks defines every gocritic checker that maps to a vint rule,
// along with its gocritic tags (sourced from go-critic.com).
var allChecks = []gocriticCheck{
	// correctness group
	{"badLock", "lint/correctness/noBadLockPattern", []string{"diagnostic", "experimental"}},
	{"badRegexp", "lint/correctness/noBadRegexpPattern", []string{"diagnostic", "experimental"}},
	{"badSorting", "lint/correctness/noBadSortUsage", []string{"diagnostic", "experimental"}},
	{"dupArg", "lint/correctness/noDuplicateArgument", []string{"diagnostic"}},
	{"dupCase", "lint/correctness/noDuplicateCase", []string{"diagnostic"}},
	{"exitAfterDefer", "lint/correctness/noExitAfterDefer", []string{"diagnostic"}},
	{"externalErrorReassign", "lint/correctness/noExternalErrorReassign", []string{"diagnostic", "experimental"}},
	{"flagDeref", "lint/correctness/noFlagDerefBeforeParse", []string{"diagnostic"}},
	{"sqlQuery", "lint/correctness/noIgnoredQueryResult", []string{"diagnostic", "experimental"}},
	{"badCond", "lint/correctness/noImpossibleCondition", []string{"diagnostic"}},
	{"badSyncOnceFunc", "lint/correctness/noInlineSyncOnceFunc", []string{"diagnostic", "experimental"}},
	{"flagName", "lint/correctness/noInvalidFlagName", []string{"diagnostic"}},
	{"returnAfterHttpError", "lint/correctness/noMissingReturnAfterHttpError", []string{"diagnostic", "experimental"}},
	{"badCall", "lint/correctness/noNoopFunctionCall", []string{"diagnostic"}},
	{"offBy1", "lint/correctness/noOffByOneError", []string{"diagnostic"}},
	{"rangeAppendAll", "lint/correctness/noRangeAppendAll", []string{"diagnostic", "experimental"}},
	{"caseOrder", "lint/correctness/noUnreachableTypeCase", []string{"diagnostic"}},
	{"preferFilepathJoin", "lint/correctness/useFilepathJoin", []string{"style", "experimental"}},

	// suspicious group
	{"sloppyLen", "lint/suspicious/noAlwaysTrueLenCheck", []string{"diagnostic"}},
	{"builtinShadow", "lint/suspicious/noBuiltinShadow", []string{"style", "opinionated"}},
	{"builtinShadowDecl", "lint/suspicious/noBuiltinShadowDecl", []string{"diagnostic", "experimental"}},
	{"deferInLoop", "lint/suspicious/noDeferInLoop", []string{"diagnostic", "experimental"}},
	{"dupBranchBody", "lint/suspicious/noDuplicateBranchBody", []string{"diagnostic"}},
	{"dupSubExpr", "lint/suspicious/noDuplicateSubExpression", []string{"diagnostic"}},
	{"dynamicFmtString", "lint/suspicious/noDynamicFormatString", []string{"diagnostic", "experimental"}},
	{"evalOrder", "lint/suspicious/noEvalOrderDependency", []string{"diagnostic", "experimental"}},
	{"importShadow", "lint/suspicious/noImportShadow", []string{"style", "opinionated"}},
	{"appendAssign", "lint/suspicious/noMismatchedAppendAssign", []string{"diagnostic"}},
	{"nilValReturn", "lint/suspicious/noNilVariableReturn", []string{"diagnostic", "experimental"}},
	{"mapKey", "lint/suspicious/noSuspiciousMapKey", []string{"diagnostic"}},
	{"sortSlice", "lint/suspicious/noSuspiciousSortSlice", []string{"diagnostic", "experimental"}},
	{"argOrder", "lint/suspicious/noSwappedArguments", []string{"diagnostic"}},
	{"truncateCmp", "lint/suspicious/noTruncatingComparison", []string{"diagnostic", "experimental"}},
	{"uncheckedInlineErr", "lint/suspicious/noUncheckedInlineError", []string{"diagnostic", "experimental"}},
	{"regexpPattern", "lint/suspicious/noUnescapedRegexpDot", []string{"diagnostic", "experimental"}},
	{"weakCond", "lint/suspicious/noWeakSliceGuard", []string{"diagnostic", "experimental"}},
	{"zeroByteRepeat", "lint/suspicious/noZeroBytesRepeat", []string{"performance"}},

	// style group
	{"captLocal", "lint/style/noCapitalizedLocal", []string{"style"}},
	{"commentedOutCode", "lint/style/noCommentedOutCode", []string{"diagnostic", "experimental"}},
	{"commentedOutImport", "lint/style/noCommentedOutImport", []string{"style", "experimental"}},
	{"docStub", "lint/style/noDocCommentStub", []string{"style", "experimental"}},
	{"dupImport", "lint/style/noDuplicateImport", []string{"style", "experimental"}},
	{"emptyDecl", "lint/style/noEmptyDeclaration", []string{"diagnostic", "experimental"}},
	{"emptyFallthrough", "lint/style/noEmptyFallthrough", []string{"style", "experimental"}},
	{"exposedSyncMutex", "lint/style/noExposedSyncMutex", []string{"style", "experimental"}},
	{"newDeref", "lint/style/noImmediateNewDeref", []string{"style"}},
	{"defaultCaseOrder", "lint/style/noMisplacedDefaultCase", []string{"style"}},
	{"hexLiteral", "lint/style/noMixedCaseHexLiteral", []string{"style", "experimental"}},
	{"ptrToRefParam", "lint/style/noPointerToRefParam", []string{"style", "opinionated", "experimental"}},
	{"unlabelStmt", "lint/style/noRedundantLabel", []string{"style", "experimental"}},
	{"unslice", "lint/style/noRedundantSliceExpression", []string{"style"}},
	{"redundantSprint", "lint/style/noRedundantSprint", []string{"style", "experimental"}},
	{"stringConcatSimplify", "lint/style/noRedundantStringConcat", []string{"style", "experimental"}},
	{"switchTrue", "lint/style/noRedundantSwitchTrue", []string{"style"}},
	{"sloppyTypeAssert", "lint/style/noRedundantTypeAssertion", []string{"diagnostic"}},
	{"filepathJoin", "lint/style/noSeparatorInFilepathJoin", []string{"diagnostic", "experimental"}},
	{"initClause", "lint/style/noSideEffectInInitClause", []string{"style", "opinionated", "experimental"}},
	{"singleCaseSwitch", "lint/style/noSingleCaseSwitch", []string{"style"}},
	{"stringsCompare", "lint/style/noStringsCompare", []string{"style", "experimental"}},
	{"todoCommentWithoutDetail", "lint/style/noTodoWithoutDetail", []string{"style", "opinionated", "experimental"}},
	{"unnecessaryBlock", "lint/style/noUnnecessaryBlock", []string{"style", "opinionated", "experimental"}},
	{"unnecessaryDefer", "lint/style/noUnnecessaryDefer", []string{"diagnostic", "experimental"}},
	{"deferUnlambda", "lint/style/noUnnecessaryDeferLambda", []string{"style", "experimental"}},
	{"underef", "lint/style/noUnnecessaryDeref", []string{"style"}},
	{"unlambda", "lint/style/noUnnecessaryLambda", []string{"style"}},
	{"typeUnparen", "lint/style/noUnnecessaryTypeParens", []string{"style", "opinionated"}},
	{"yodaStyleExpr", "lint/style/noYodaCondition", []string{"style", "experimental"}},
	{"assignOp", "lint/style/useAssignmentOperator", []string{"style"}},
	{"paramTypeCombine", "lint/style/useCombinedParamType", []string{"style", "opinionated"}},
	{"commentFormatting", "lint/style/useCommentSpacing", []string{"style"}},
	{"wrapperFunc", "lint/style/useConvenienceFunc", []string{"style"}},
	{"methodExprCall", "lint/style/useDirectMethodCall", []string{"style", "experimental"}},
	{"emptyStringTest", "lint/style/useDirectStringComparison", []string{"style", "experimental"}},
	{"elseif", "lint/style/useElseIf", []string{"style"}},
	{"httpNoBody", "lint/style/useHttpNoBody", []string{"style", "experimental"}},
	{"octalLiteral", "lint/style/useModernOctalLiteral", []string{"style", "experimental", "opinionated"}},
	{"unnamedResult", "lint/style/useNamedResult", []string{"style", "opinionated", "experimental"}},
	{"valSwap", "lint/style/useParallelAssignSwap", []string{"style"}},
	{"sprintfQuotedString", "lint/style/usePercentQ", []string{"diagnostic", "experimental"}},
	{"regexpMust", "lint/style/useRegexpMustCompile", []string{"style"}},
	{"sloppyReassign", "lint/style/useShortVarDecl", []string{"diagnostic", "experimental"}},
	{"regexpSimplify", "lint/style/useSimplifiedRegexp", []string{"style", "experimental", "opinionated"}},
	{"codegenComment", "lint/style/useStandardCodegenComment", []string{"diagnostic"}},
	{"deprecatedComment", "lint/style/useStandardDeprecationComment", []string{"diagnostic"}},
	{"ifElseChain", "lint/style/useSwitch", []string{"style"}},
	{"timeExprSimplify", "lint/style/useTimeMethod", []string{"style", "experimental"}},
	{"typeDefFirst", "lint/style/useTypeDefFirst", []string{"style", "experimental"}},
	{"typeAssertChain", "lint/style/useTypeSwitchChain", []string{"style", "experimental"}},
	{"typeSwitchVar", "lint/style/useTypeSwitchGuard", []string{"style"}},

	// performance group
	{"hugeParam", "lint/performance/noHugeParam", []string{"performance"}},
	{"rangeExprCopy", "lint/performance/noRangeExprCopy", []string{"performance"}},
	{"rangeValCopy", "lint/performance/noRangeValCopy", []string{"performance"}},
	{"stringXbytes", "lint/performance/noRedundantStringByteConversion", []string{"performance"}},
	{"indexAlloc", "lint/performance/noStringIndexAllocation", []string{"performance"}},
	{"appendCombine", "lint/performance/useCombinedAppend", []string{"performance"}},
	{"preferDecodeRune", "lint/performance/useDecodeRune", []string{"performance", "experimental"}},
	{"equalFold", "lint/performance/useEqualFold", []string{"performance", "experimental"}},
	{"preferFprint", "lint/performance/useFprint", []string{"performance", "experimental"}},
	{"sliceClear", "lint/performance/useOptimizedSliceClear", []string{"performance", "experimental"}},
	{"preferStringWriter", "lint/performance/useStringWriter", []string{"performance", "experimental"}},
	{"syncMapLoadAndDelete", "lint/performance/useSyncMapLoadAndDelete", []string{"diagnostic", "experimental"}},
	{"preferWriteByte", "lint/performance/useWriteByte", []string{"performance", "experimental", "opinionated"}},

	// complexity group
	{"tooManyResultsChecker", "lint/complexity/noExcessiveResults", []string{"style", "opinionated", "experimental"}},
	{"nestingReduce", "lint/complexity/useEarlyContinue", []string{"style", "opinionated", "experimental"}},
	{"boolExprSimplify", "lint/complexity/useSimplifiedBoolExpr", []string{"style", "experimental"}},
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	enabled := computeEnabledCheckers(settings)

	configs := make(map[string]migrate.VintRuleConfig)
	for _, check := range allChecks {
		if enabled[check.checker] {
			configs[check.vintRule] = migrate.VintRuleConfig{}
		}
	}

	// Map gocritic-specific settings to individual rules when present.
	if settings != nil {
		if settingsMap, ok := settings["settings"]; ok {
			if sm, ok := settingsMap.(map[string]any); ok {
				applyCheckerSettings(configs, sm)
			}
		}
	}

	return configs, nil
}

// applyCheckerSettings maps gocritic per-checker options to the corresponding
// vint rules. Settings are only applied if the rule is already enabled.
func applyCheckerSettings(configs map[string]migrate.VintRuleConfig, sm map[string]any) {
	// hugeParam sizeThreshold
	if hugeParam, ok := sm["hugeParam"]; ok {
		if hp, ok := hugeParam.(map[string]any); ok {
			if threshold, ok := hp["sizeThreshold"]; ok {
				setOptionsIfEnabled(configs, "lint/performance/noHugeParam", map[string]any{
					"sizeThreshold": threshold,
				})
			}
		}
	}

	// rangeValCopy sizeThreshold and skipTestFuncs
	if rangeValCopy, ok := sm["rangeValCopy"]; ok {
		if rvc, ok := rangeValCopy.(map[string]any); ok {
			opts := map[string]any{}
			if threshold, ok := rvc["sizeThreshold"]; ok {
				opts["sizeThreshold"] = threshold
			}
			if skipTests, ok := rvc["skipTestFuncs"]; ok {
				opts["skipTestFuncs"] = skipTests
			}
			if len(opts) > 0 {
				setOptionsIfEnabled(configs, "lint/performance/noRangeValCopy", opts)
			}
		}
	}

	// rangeExprCopy sizeThreshold and skipTestFuncs
	if rangeExprCopy, ok := sm["rangeExprCopy"]; ok {
		if rec, ok := rangeExprCopy.(map[string]any); ok {
			opts := map[string]any{}
			if threshold, ok := rec["sizeThreshold"]; ok {
				opts["sizeThreshold"] = threshold
			}
			if skipTests, ok := rec["skipTestFuncs"]; ok {
				opts["skipTestFuncs"] = skipTests
			}
			if len(opts) > 0 {
				setOptionsIfEnabled(configs, "lint/performance/noRangeExprCopy", opts)
			}
		}
	}

	// tooManyResultsChecker -> noExcessiveResults maxResults
	if tooManyResults, ok := sm["tooManyResultsChecker"]; ok {
		if tmr, ok := tooManyResults.(map[string]any); ok {
			if maxResults, ok := tmr["maxResults"]; ok {
				setOptionsIfEnabled(configs, "lint/complexity/noExcessiveResults", map[string]any{
					"maxResults": maxResults,
				})
			}
		}
	}

	// ifElseChain -> useSwitch minCases
	if ifElseChain, ok := sm["ifElseChain"]; ok {
		if iec, ok := ifElseChain.(map[string]any); ok {
			if minCases, ok := iec["minCases"]; ok {
				setOptionsIfEnabled(configs, "lint/style/useSwitch", map[string]any{
					"minCases": minCases,
				})
			}
		}
	}

	// nestingReduce -> useEarlyContinue bodyWidth
	if nestingReduce, ok := sm["nestingReduce"]; ok {
		if nr, ok := nestingReduce.(map[string]any); ok {
			if bodyWidth, ok := nr["bodyWidth"]; ok {
				setOptionsIfEnabled(configs, "lint/complexity/useEarlyContinue", map[string]any{
					"bodyWidth": bodyWidth,
				})
			}
		}
	}

	// captLocal -> noCapitalizedLocal paramsOnly
	if captLocal, ok := sm["captLocal"]; ok {
		if cl, ok := captLocal.(map[string]any); ok {
			if paramsOnly, ok := cl["paramsOnly"]; ok {
				setOptionsIfEnabled(configs, "lint/style/noCapitalizedLocal", map[string]any{
					"paramsOnly": paramsOnly,
				})
			}
		}
	}

	// unnamedResult -> useNamedResult maxResults
	if unnamedResult, ok := sm["unnamedResult"]; ok {
		if ur, ok := unnamedResult.(map[string]any); ok {
			if maxResults, ok := ur["maxResults"]; ok {
				setOptionsIfEnabled(configs, "lint/style/useNamedResult", map[string]any{
					"maxResults": maxResults,
				})
			}
		}
	}
}

// setOptionsIfEnabled updates a rule's options only if the rule is already present
// in the configs map (i.e., the corresponding gocritic check is enabled).
func setOptionsIfEnabled(configs map[string]migrate.VintRuleConfig, rule string, opts map[string]any) {
	if _, ok := configs[rule]; ok {
		configs[rule] = migrate.VintRuleConfig{Options: opts}
	}
}

// computeEnabledCheckers determines which gocritic checks are enabled
// based on the golangci-lint settings. The algorithm:
//  1. Start with the default set (non-experimental, non-opinionated checks)
//  2. Add all checks matching any tag in enabled-tags
//  3. Remove all checks matching any tag in disabled-tags
//  4. Add specific checks from enabled-checks
//  5. Remove specific checks from disabled-checks
func computeEnabledCheckers(settings map[string]any) map[string]bool {
	enabled := make(map[string]bool)

	// Step 1: defaults — enable all non-experimental, non-opinionated checks.
	for _, check := range allChecks {
		if isDefaultEnabled(check.tags) {
			enabled[check.checker] = true
		}
	}

	if settings == nil {
		return enabled
	}

	// Step 2: enabled-tags — add all checks matching any enabled tag.
	if enabledTags := getStringSlice(settings, "enabled-tags"); len(enabledTags) > 0 {
		tagSet := toSet(enabledTags)
		for _, check := range allChecks {
			if hasAnyTag(check.tags, tagSet) {
				enabled[check.checker] = true
			}
		}
	}

	// Step 3: disabled-tags — remove all checks matching any disabled tag.
	if disabledTags := getStringSlice(settings, "disabled-tags"); len(disabledTags) > 0 {
		tagSet := toSet(disabledTags)
		for _, check := range allChecks {
			if hasAnyTag(check.tags, tagSet) {
				delete(enabled, check.checker)
			}
		}
	}

	// Step 4: enabled-checks — add specific checks.
	for _, name := range getStringSlice(settings, "enabled-checks") {
		enabled[name] = true
	}

	// Step 5: disabled-checks — remove specific checks.
	for _, name := range getStringSlice(settings, "disabled-checks") {
		delete(enabled, name)
	}

	return enabled
}

// isDefaultEnabled returns true if a checker's tags indicate it should be
// enabled by default (not experimental and not opinionated).
func isDefaultEnabled(tags []string) bool {
	for _, tag := range tags {
		if tag == "experimental" || tag == "opinionated" {
			return false
		}
	}
	return true
}

// hasAnyTag returns true if any of the checker's tags is in the given set.
func hasAnyTag(tags []string, tagSet map[string]bool) bool {
	for _, tag := range tags {
		if tagSet[tag] {
			return true
		}
	}
	return false
}

// getStringSlice extracts a []string from a settings map entry.
// YAML unmarshalling may produce []any instead of []string.
func getStringSlice(settings map[string]any, key string) []string {
	v, ok := settings[key]
	if !ok {
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
	}
	return nil
}

// toSet converts a string slice to a set for efficient lookup.
func toSet(s []string) map[string]bool {
	m := make(map[string]bool, len(s))
	for _, v := range s {
		m[v] = true
	}
	return m
}

package gocritic

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the gocritic golangci-lint linter.
// gocritic is a large linter with many checks covering correctness,
// performance, style, and complexity. All checks are enabled by default
// when gocritic is enabled.
type Migrator struct{}

func (*Migrator) Name() string {
	return "gocritic"
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// gocritic enables all its checks by default, so all mapped rules
	// are unconditionally enabled.

	// correctness group
	configs["lint/correctness/noBadLockPattern"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noBadRegexpPattern"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noBadSortUsage"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noDuplicateArgument"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noDuplicateCase"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noExitAfterDefer"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noExternalErrorReassign"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noFlagDerefBeforeParse"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noIgnoredQueryResult"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noImpossibleCondition"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noInlineSyncOnceFunc"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noInvalidFlagName"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noMissingReturnAfterHttpError"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noNoopFunctionCall"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noOffByOneError"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noRangeAppendAll"] = migrate.VintRuleConfig{}
	configs["lint/correctness/noUnreachableTypeCase"] = migrate.VintRuleConfig{}
	configs["lint/correctness/useFilepathJoin"] = migrate.VintRuleConfig{}

	// suspicious group
	configs["lint/suspicious/noAlwaysTrueLenCheck"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noBuiltinShadow"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noBuiltinShadowDecl"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noDeferInLoop"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noDuplicateBranchBody"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noDuplicateSubExpression"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noDynamicFormatString"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noEvalOrderDependency"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noImportShadow"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noMismatchedAppendAssign"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noNilVariableReturn"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noSuspiciousMapKey"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noSuspiciousSortSlice"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noSwappedArguments"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noTruncatingComparison"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noUncheckedInlineError"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noUnescapedRegexpDot"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noWeakSliceGuard"] = migrate.VintRuleConfig{}
	configs["lint/suspicious/noZeroBytesRepeat"] = migrate.VintRuleConfig{}

	// style group
	configs["lint/style/noCapitalizedLocal"] = migrate.VintRuleConfig{}
	configs["lint/style/noCommentedOutCode"] = migrate.VintRuleConfig{}
	configs["lint/style/noCommentedOutImport"] = migrate.VintRuleConfig{}
	configs["lint/style/noDocCommentStub"] = migrate.VintRuleConfig{}
	configs["lint/style/noDuplicateImport"] = migrate.VintRuleConfig{}
	configs["lint/style/noEmptyDeclaration"] = migrate.VintRuleConfig{}
	configs["lint/style/noEmptyFallthrough"] = migrate.VintRuleConfig{}
	configs["lint/style/noExposedSyncMutex"] = migrate.VintRuleConfig{}
	configs["lint/style/noImmediateNewDeref"] = migrate.VintRuleConfig{}
	configs["lint/style/noMisplacedDefaultCase"] = migrate.VintRuleConfig{}
	configs["lint/style/noMixedCaseHexLiteral"] = migrate.VintRuleConfig{}
	configs["lint/style/noPointerToRefParam"] = migrate.VintRuleConfig{}
	configs["lint/style/noRedundantLabel"] = migrate.VintRuleConfig{}
	configs["lint/style/noRedundantSliceExpression"] = migrate.VintRuleConfig{}
	configs["lint/style/noRedundantSprint"] = migrate.VintRuleConfig{}
	configs["lint/style/noRedundantStringConcat"] = migrate.VintRuleConfig{}
	configs["lint/style/noRedundantSwitchTrue"] = migrate.VintRuleConfig{}
	configs["lint/style/noRedundantTypeAssertion"] = migrate.VintRuleConfig{}
	configs["lint/style/noSeparatorInFilepathJoin"] = migrate.VintRuleConfig{}
	configs["lint/style/noSideEffectInInitClause"] = migrate.VintRuleConfig{}
	configs["lint/style/noSingleCaseSwitch"] = migrate.VintRuleConfig{}
	configs["lint/style/noStringsCompare"] = migrate.VintRuleConfig{}
	configs["lint/style/noTodoWithoutDetail"] = migrate.VintRuleConfig{}
	configs["lint/style/noUnnecessaryBlock"] = migrate.VintRuleConfig{}
	configs["lint/style/noUnnecessaryDefer"] = migrate.VintRuleConfig{}
	configs["lint/style/noUnnecessaryDeferLambda"] = migrate.VintRuleConfig{}
	configs["lint/style/noUnnecessaryDeref"] = migrate.VintRuleConfig{}
	configs["lint/style/noUnnecessaryLambda"] = migrate.VintRuleConfig{}
	configs["lint/style/noUnnecessaryTypeParens"] = migrate.VintRuleConfig{}
	configs["lint/style/noYodaCondition"] = migrate.VintRuleConfig{}
	configs["lint/style/useAssignmentOperator"] = migrate.VintRuleConfig{}
	configs["lint/style/useCombinedParamType"] = migrate.VintRuleConfig{}
	configs["lint/style/useCommentSpacing"] = migrate.VintRuleConfig{}
	configs["lint/style/useConvenienceFunc"] = migrate.VintRuleConfig{}
	configs["lint/style/useDirectMethodCall"] = migrate.VintRuleConfig{}
	configs["lint/style/useDirectStringComparison"] = migrate.VintRuleConfig{}
	configs["lint/style/useElseIf"] = migrate.VintRuleConfig{}
	configs["lint/style/useHttpNoBody"] = migrate.VintRuleConfig{}
	configs["lint/style/useModernOctalLiteral"] = migrate.VintRuleConfig{}
	configs["lint/style/useNamedResult"] = migrate.VintRuleConfig{}
	configs["lint/style/useParallelAssignSwap"] = migrate.VintRuleConfig{}
	configs["lint/style/usePercentQ"] = migrate.VintRuleConfig{}
	configs["lint/style/useRegexpMustCompile"] = migrate.VintRuleConfig{}
	configs["lint/style/useShortVarDecl"] = migrate.VintRuleConfig{}
	configs["lint/style/useSimplifiedRegexp"] = migrate.VintRuleConfig{}
	configs["lint/style/useStandardCodegenComment"] = migrate.VintRuleConfig{}
	configs["lint/style/useStandardDeprecationComment"] = migrate.VintRuleConfig{}
	configs["lint/style/useSwitch"] = migrate.VintRuleConfig{}
	configs["lint/style/useTimeMethod"] = migrate.VintRuleConfig{}
	configs["lint/style/useTypeDefFirst"] = migrate.VintRuleConfig{}
	configs["lint/style/useTypeSwitchChain"] = migrate.VintRuleConfig{}
	configs["lint/style/useTypeSwitchGuard"] = migrate.VintRuleConfig{}

	// performance group
	configs["lint/performance/noHugeParam"] = migrate.VintRuleConfig{}
	configs["lint/performance/noRangeExprCopy"] = migrate.VintRuleConfig{}
	configs["lint/performance/noRangeValCopy"] = migrate.VintRuleConfig{}
	configs["lint/performance/noRedundantStringByteConversion"] = migrate.VintRuleConfig{}
	configs["lint/performance/noStringIndexAllocation"] = migrate.VintRuleConfig{}
	configs["lint/performance/useCombinedAppend"] = migrate.VintRuleConfig{}
	configs["lint/performance/useDecodeRune"] = migrate.VintRuleConfig{}
	configs["lint/performance/useEqualFold"] = migrate.VintRuleConfig{}
	configs["lint/performance/useFprint"] = migrate.VintRuleConfig{}
	configs["lint/performance/useOptimizedSliceClear"] = migrate.VintRuleConfig{}
	configs["lint/performance/useStringWriter"] = migrate.VintRuleConfig{}
	configs["lint/performance/useSyncMapLoadAndDelete"] = migrate.VintRuleConfig{}
	configs["lint/performance/useWriteByte"] = migrate.VintRuleConfig{}

	// complexity group
	configs["lint/complexity/noExcessiveResults"] = migrate.VintRuleConfig{}
	configs["lint/complexity/useEarlyContinue"] = migrate.VintRuleConfig{}
	configs["lint/complexity/useSimplifiedBoolExpr"] = migrate.VintRuleConfig{}

	// Map gocritic-specific settings to individual rules when present.
	if settings != nil {
		// hugeParam sizeThreshold
		if settingsMap, ok := settings["settings"]; ok {
			if sm, ok := settingsMap.(map[string]any); ok {
				if hugeParam, ok := sm["hugeParam"]; ok {
					if hp, ok := hugeParam.(map[string]any); ok {
						if threshold, ok := hp["sizeThreshold"]; ok {
							configs["lint/performance/noHugeParam"] = migrate.VintRuleConfig{
								Options: map[string]any{"sizeThreshold": threshold},
							}
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
							configs["lint/performance/noRangeValCopy"] = migrate.VintRuleConfig{
								Options: opts,
							}
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
							configs["lint/performance/noRangeExprCopy"] = migrate.VintRuleConfig{
								Options: opts,
							}
						}
					}
				}

				// tooManyResults -> noExcessiveResults maxResults
				if tooManyResults, ok := sm["tooManyResultsChecker"]; ok {
					if tmr, ok := tooManyResults.(map[string]any); ok {
						if maxResults, ok := tmr["maxResults"]; ok {
							configs["lint/complexity/noExcessiveResults"] = migrate.VintRuleConfig{
								Options: map[string]any{"maxResults": maxResults},
							}
						}
					}
				}

				// ifElseChain -> useSwitch minCases
				if ifElseChain, ok := sm["ifElseChain"]; ok {
					if iec, ok := ifElseChain.(map[string]any); ok {
						if minCases, ok := iec["minCases"]; ok {
							configs["lint/style/useSwitch"] = migrate.VintRuleConfig{
								Options: map[string]any{"minCases": minCases},
							}
						}
					}
				}

				// nestingReduce -> useEarlyContinue bodyWidth
				if nestingReduce, ok := sm["nestingReduce"]; ok {
					if nr, ok := nestingReduce.(map[string]any); ok {
						if bodyWidth, ok := nr["bodyWidth"]; ok {
							configs["lint/complexity/useEarlyContinue"] = migrate.VintRuleConfig{
								Options: map[string]any{"bodyWidth": bodyWidth},
							}
						}
					}
				}

				// captLocal -> noCapitalizedLocal paramsOnly
				if captLocal, ok := sm["captLocal"]; ok {
					if cl, ok := captLocal.(map[string]any); ok {
						if paramsOnly, ok := cl["paramsOnly"]; ok {
							configs["lint/style/noCapitalizedLocal"] = migrate.VintRuleConfig{
								Options: map[string]any{"paramsOnly": paramsOnly},
							}
						}
					}
				}

				// unnamedResult -> useNamedResult maxResults
				if unnamedResult, ok := sm["unnamedResult"]; ok {
					if ur, ok := unnamedResult.(map[string]any); ok {
						if maxResults, ok := ur["maxResults"]; ok {
							configs["lint/style/useNamedResult"] = migrate.VintRuleConfig{
								Options: map[string]any{"maxResults": maxResults},
							}
						}
					}
				}
			}
		}
	}

	return configs, nil
}

package config_test

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	goversion "github.com/hashicorp/go-version"

	"github.com/strowk/vint/config"
	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/no_deep_exit"
	"github.com/strowk/vint/rules/no_excessive_arguments"
)

func TestGetConfig(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		for name, tc := range map[string]struct {
			confPath   string
			wantConfig lint.Config
		}{
			"default config": {
				wantConfig: lint.Config{
					IgnoreGeneratedHeader: false,
					Confidence:            0.8,
					Severity:              lint.SeverityWarning,
					EnableAllRules:        false,
					EnableDefaultRules:    false,
					Rules: lint.RulesConfig{
						"lint/style/noBlankImport": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useContextAsFirstParam": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noContextKeysType": {
							Severity: lint.SeverityWarning,
						},
						"empty-block": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useErrorNaming": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noErrorStrings": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useErrorf": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useExportedComment": {
							Severity: lint.SeverityWarning,
						},
						"increment-decrement": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useIndentErrorFlow": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnclosedBodies": {
							Severity: lint.SeverityWarning,
						},
						"lint/complexity/noDuplicateCode": {
							Severity: lint.SeverityWarning,
						},
						"lint/complexity/noDuplicateCodeV2": {
							Severity: lint.SeverityWarning,
						},
						"lint/complexity/noExcessiveStatements": {
							Severity: lint.SeverityWarning,
						},
						"lint/complexity/noHighCyclomaticComplexity": {
							Severity: lint.SeverityWarning,
						},
						"lint/complexity/noLongFunctions": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noDeniedImport": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noDeprecatedUsage": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noDirectErrorComparison": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noDiscardedAppend": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noDeferCloseBeforeErrCheck": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noDeferInInfiniteLoop": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noDynamicErrors": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noEmptyForLoop": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSelectBreakConfusion": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noInfiniteRecursion": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noFileScopedDeniedImport": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSpaceInDirective": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnallowedImport": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUncheckedError": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noRangeVariableAlias": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSliceBoundsOutOfRange": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnrecognizedDirective": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noAtoiOverflow": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noBindToAllInterfaces": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noExposedPprof": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noFilesystemRootServing": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noHardcodedCredentials": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noHardcodedIv": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noHttpRequestSmuggling": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noInsecureCookie": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noInsecureHostKeyCallback": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noIntegerOverflowConversion": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noMissingReadHeaderTimeout": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noServeWithoutTimeout": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noSsrfViaVariable": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noTrojanSourceBidi": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noUnboundedDecompression": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noUnsafePackage": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noCgiImport": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noCommandInjectionTaint": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noContextPropagationFailure": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noDeprecatedHashFunction": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noFileInclusionViaVariable": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noFilesystemToctou": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noInsecureRandom": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noInsecureTlsConfig": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noLogInjectionTaint": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noPathTraversalTaint": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noPermissiveDirectoryPermissions": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noPermissiveFilePermissions": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noPermissiveOsCreate": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noPermissiveWriteFilePermissions": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noPredictableTempFile": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noSecretInSerialization": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noShortRsaKey": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noSmtpInjectionTaint": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noSqlConcatenation": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noSqlFormatString": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noSqlInjectionTaint": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noSshAuthBypass": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noSsrfTaint": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noTemplateInjection": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noTlsSessionResumptionBypass": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noUnboundedFormParsing": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noUnsafeCorsBypass": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noUnsafeRedirectPolicy": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noUnescapedHtmlTemplate": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noUnsafeDeserialization": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noVariableCommandExecution": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noWeakCryptoHash": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noWeakEncryptionAlgorithm": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noXssTaint": {
							Severity: lint.SeverityWarning,
						},
						"lint/security/noZipSlip": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noExcessiveBlankIdentifiers": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noDotImport": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noInitFunction": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noLeadingBlankLine": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noTrailingBlankLine": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noRepeatedStrings": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noUnnecessaryConversion": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noUnnecessaryDeref": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noUnnecessaryBlankIdentifier": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noUnnecessaryLoopVarCopy": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/usePrintfSuffix": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSingleArgAppend": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noAsmDeclMismatch": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSelfAssignment": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSelfReferencingFinalizer": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noAtomicAssignMisuse": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noMalformedBuildTag": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noMalformedDirective": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noMalformedStructTag": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noCgoPointerViolation": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noUnkeyedLiteral": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noCopiedLock": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noDeferTimeMisuse": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noIncorrectTimeFormat": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noHttpClientGet": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noHttpClientHead": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noHttpClientPost": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noHttpClientPostForm": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noHttpGet": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noHttpHead": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noHttpPost": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noHttpPostForm": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noHttpNewRequest": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noHttpResponseMisuse": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noHttptestNewRequest": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noInvalidBinaryArg": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noInvalidErrorsAs": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noFramePointerClobber": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/useJoinHostPort": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noImpossibleInterfaceAssert": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noLoopClosureCapture": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noLostCancel": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noWaitGroupMisuse": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noWriterBufferModification": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNilContext": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNilFuncComparison": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNilMapAssignment": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noPrintfFormatMismatch": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noExcessiveShift": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnbufferedSignalChannel": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUntrappableSignal": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSlogKeyValueMismatch": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noOddSizeSliceArg": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noStdMethodSignatureMismatch": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noStdlibVersionMismatch": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noStringIntConversion": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noTestFatalInGoroutine": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noTestMainWithoutExit": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noMalformedTestFunction": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNonCanonicalHeaderKey": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noRedundantCanonicalHeaderKey": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNonPointerUnmarshal": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnmarshalableType": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnreachableCode": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noInvalidUnsafePointer": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnusedFunctionResult": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnusedWrite": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noIneffectualAssignment": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noIneffectiveBreak": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noLineTooLong": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noMagicNumberInArgument": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noMagicNumberInAssignment": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noMagicNumberInCase": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noMagicNumberInCondition": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noMagicNumberInOperation": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noMagicNumberInReturn": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noMisspelledWords": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNanComparison": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noNakedReturn": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noTlsConnHandshake": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noTlsDial": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noTlsDialWithDialer": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNetDial": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNetDialTimeout": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNetIpBytesEqual": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNetListen": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNetListenPacket": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNetLookupCname": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNetLookupHost": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNetLookupIp": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNetLookupMx": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNetLookupPort": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNetLookupSrv": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNetLookupNs": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNetLookupTxt": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNetLookupAddr": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSqlDbBegin": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSqlDbExec": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSqlDbPing": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSqlDbPrepare": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSqlDbQuery": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSqlDbQueryRow": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSqlTxExec": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSqlTxPrepare": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSqlTxQuery": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSqlTxQueryRow": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSqlTxStmt": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSqlStmtExec": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSqlStmtQuery": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noSqlStmtQueryRow": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noExecCommand": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noInvalidExecCommandArg": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noNolintNonMachineReadable": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNolintParseError": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noNolintWithoutSpecificLinter": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noNolintWithoutExplanation": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noUnusedNolint": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noUnusedParameter": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noConstantParameter": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noConstantResult": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnusedFunction": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnusedType": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnusedVariable": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnusedConstant": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnusedField": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noMismatchedAppendAssign": {
							Severity: lint.SeverityWarning,
						},
						"lint/performance/noHugeParam": {
							Severity: lint.SeverityWarning,
						},
						"lint/performance/noStringIndexAllocation": {
							Severity: lint.SeverityWarning,
						},
						"lint/performance/useStringMapKey": {
							Severity: lint.SeverityWarning,
						},
						"lint/performance/noRedundantStringByteConversion": {
							Severity: lint.SeverityWarning,
						},
						"lint/performance/useInlineMathPow": {
							Severity: lint.SeverityWarning,
						},
						"lint/performance/useCombinedAppend": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noSwappedArguments": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useAssignmentOperator": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNoopFunctionCall": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noImpossibleCondition": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noCapitalizedLocal": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noCapitalizedErrorString": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnreachableTypeCase": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noMisplacedDefaultCase": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useStandardCodegenComment": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useStandardDeprecationComment": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useCommentSpacing": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noDuplicateBuildConstraint": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noDuplicateArgument": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noDuplicateBranchBody": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noDuplicateCase": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noDuplicateSubExpression": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noDuplicateCutsetChars": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noDuplicateIfCondition": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useElseIf": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noExitAfterDefer": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noFlagDerefBeforeParse": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noInvalidFlagName": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noInvalidHostPort": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useSwitch": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useTaggedSwitch": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useSimplifiedBoolReturn": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useSimplifiedPrintFormat": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useStringConversionInPrint": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/usePackageComment": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useIdiomaticErrorName": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useIdiomaticNaming": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useIdiomaticReceiverName": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useIdiomaticDurationName": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useConsistentReceiverName": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noControlCharInString": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useFuncDocPrefix": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useTypeDocPrefix": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useVarConstDocPrefix": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noRedundantMakeArgs": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noRedundantNilLoopCheck": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noRedundantNilSliceCheck": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noRedundantNilTypeCheck": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noDefaultSliceIndex": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useSliceAppend": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/useTimeEqual": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useTimeSince": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useTimeSleep": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useTimeUntil": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useTypeConversion": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useTrimFunction": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useMergedConditionalDecl": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useMergedVarDecl": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noRedundantControlFlow": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useErrorMethod": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useErrorLastReturn": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useDirectStringRange": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useBufferStringOrBytes": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noGuardAroundDelete": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noGuardAroundMapAccess": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noSuspiciousMapKey": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noSuspiciousTimeSleep": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noEmptyBranch": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noEmptyCriticalSection": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noTimeTick": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noTimerResetRetval": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noImmediateNewDeref": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noOffByOneError": {
							Severity: lint.SeverityWarning,
						},
						"lint/performance/noRangeExprCopy": {
							Severity: lint.SeverityWarning,
						},
						"lint/performance/noRangeValCopy": {
							Severity: lint.SeverityWarning,
						},
						"lint/performance/noRedundantRuneConversion": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useRegexpMustCompile": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useRawStringRegexp": {
							Severity: lint.SeverityWarning,
						},
						"lint/performance/useCompiledRegexp": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noSingleCaseSelect": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noSingleCaseSwitch": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noRedundantSwitchTrue": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noAlwaysTrueLenCheck": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noRedundantTypeAssertion": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useTypeSwitchGuard": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useTypeAssertResult": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noUnnecessaryLambda": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noRedundantSliceExpression": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useParallelAssignSwap": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useConvenienceFunc": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noZeroBytesRepeat": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noInvalidRegexp": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noInvalidTemplate": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noInvalidUrlParse": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noInvalidUtf8StringArg": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noInappropriateContextKey": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noInvalidStrconvArg": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noOverlappingEncoderSlice": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noBenchmarkNAssignment": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noAddressOfDereference": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noAddressNilComparison": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noUnsignedNegativeComparison": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noUnobservedFieldAssign": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noOverwrittenArgument": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noSingleIterationLoop": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noInvariantLoopCondition": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noUselessMathCall": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noIneffectiveBitwiseOp": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noIneffectiveRandCall": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noImpossibleNilComparison": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noImpossibleBuiltinResult": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noIntegerDivisionTruncation": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noImpreciseConstant": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noImplicitConstValue": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noIgnoredQueryModification": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noModuloOne": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noNeverNilCheck": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noNonOctalFileMode": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noUnmarshalableStruct": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noDubiousBitShift": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noTempDirDeletion": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noTypeAssertElseMisread": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useCopyBuiltin": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useCopyForSlide": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noExplicitBoolComparison": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useStringsContains": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useBytesEqual": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useInfiniteFor": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noRedundantVarType": {
							Severity: lint.SeverityWarning,
						},
						"lint/complexity/useLoopCondition": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useSimplifiedSelector": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/usePackageComments": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noRedundantRangeVal": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useReceiverNaming": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noSuperfluousElse": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useTimeNaming": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noUnexportedReturn": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/noVarDeclaration": {
							Severity: lint.SeverityWarning,
						},
						"lint/style/useVarNaming": {
							Severity: lint.SeverityWarning,
						},
						"lint/performance/usePointerInSyncPool": {
							Severity: lint.SeverityWarning,
						},
					},
					ErrorCode:   0,
					WarningCode: 0,
					Directives:  lint.DirectivesConfig{},
					Exclude:     []string{},
					GoVersion:   nil,
				},
			},
			"non-reg issue #470": {
				confPath: "issue-470.toml",
				wantConfig: lint.Config{
					Confidence: 0.8,
					Severity:   lint.SeverityWarning,
					Rules: lint.RulesConfig{
						"add-constant": {
							Severity: lint.SeverityWarning,
							Arguments: lint.Arguments{
								map[string]any{
									"maxLitCount": "3",
									"allowStrs":   `"`,
									"allowFloats": "0.0,1.0,1.,2.0,2.",
									"allowInts":   "0,1,2",
								},
							},
						},
					},
				},
			},
			"config from file issue #585": {
				confPath: "issue-585.toml",
				wantConfig: lint.Config{
					Confidence: 0.0,
					Severity:   lint.SeverityWarning,
				},
			},
			"config from file default confidence issue #585": {
				confPath: "issue-585-defaultConfidence.toml",
				wantConfig: lint.Config{
					Confidence: 0.8,
					Severity:   lint.SeverityWarning,
				},
			},
			"config from file goVersion": {
				confPath: "goVersion.toml",
				wantConfig: lint.Config{
					Confidence: 0.8,
					GoVersion:  goversion.Must(goversion.NewSemver("1.20.0")),
				},
			},
			"config from file ignoreGeneratedHeader": {
				confPath: "ignoreGeneratedHeader.toml",
				wantConfig: lint.Config{
					Confidence:            0.8,
					IgnoreGeneratedHeader: true,
				},
			},
			"config from file enableDefault": {
				confPath: "enableDefault.toml",
				wantConfig: lint.Config{
					Confidence:            0.8,
					IgnoreGeneratedHeader: false,
					EnableDefaultRules:    true,
					Rules: lint.RulesConfig{
						"lint/style/noBlankImport":              {},
						"lint/style/useContextAsFirstParam":                   {},
						"lint/correctness/noContextKeysType":                     {},
						"empty-block":                           {},
						"lint/style/useErrorNaming":                          {},
						"lint/style/noErrorStrings":                         {},
						"lint/style/useErrorf":                                {},
						"lint/style/useExportedComment":                              {},
						"increment-decrement":                   {},
						"lint/style/useIndentErrorFlow":                     {},
						"lint/correctness/noUnclosedBodies":     {},
						"lint/complexity/noDuplicateCode":       {},
						"lint/complexity/noDuplicateCodeV2":     {},
						"lint/complexity/noExcessiveStatements": {},
						"lint/complexity/noHighCyclomaticComplexity":     {},
						"lint/complexity/noLongFunctions":                {},
						"lint/correctness/noDeniedImport":                {},
						"lint/correctness/noDeprecatedUsage":             {},
						"lint/correctness/noDirectErrorComparison":       {},
						"lint/correctness/noDiscardedAppend":             {},
						"lint/correctness/noDeferCloseBeforeErrCheck":    {},
						"lint/correctness/noDeferInInfiniteLoop":         {},
						"lint/correctness/noDynamicErrors":               {},
						"lint/correctness/noEmptyForLoop":                {},
						"lint/correctness/noSelectBreakConfusion":        {},
						"lint/correctness/noInfiniteRecursion":           {},
						"lint/correctness/noFileScopedDeniedImport":      {},
						"lint/correctness/noSpaceInDirective":            {},
						"lint/correctness/noUnallowedImport":             {},
						"lint/correctness/noUncheckedError":              {},
						"lint/correctness/noRangeVariableAlias":          {},
						"lint/correctness/noSliceBoundsOutOfRange":       {},
						"lint/correctness/noUnrecognizedDirective":       {},
						"lint/security/noAtoiOverflow":                   {},
						"lint/security/noBindToAllInterfaces":            {},
						"lint/security/noExposedPprof":                   {},
						"lint/security/noFilesystemRootServing":          {},
						"lint/security/noHardcodedCredentials":           {},
						"lint/security/noHardcodedIv":                    {},
						"lint/security/noHttpRequestSmuggling":           {},
						"lint/security/noInsecureCookie":                 {},
						"lint/security/noInsecureHostKeyCallback":        {},
						"lint/security/noIntegerOverflowConversion":      {},
						"lint/security/noMissingReadHeaderTimeout":       {},
						"lint/security/noServeWithoutTimeout":            {},
						"lint/security/noSsrfViaVariable":                {},
						"lint/security/noTrojanSourceBidi":               {},
						"lint/security/noUnboundedDecompression":         {},
						"lint/security/noUnsafePackage":                  {},
						"lint/security/noCgiImport":                      {},
						"lint/security/noCommandInjectionTaint":          {},
						"lint/security/noContextPropagationFailure":      {},
						"lint/security/noDeprecatedHashFunction":         {},
						"lint/security/noFileInclusionViaVariable":       {},
						"lint/security/noFilesystemToctou":               {},
						"lint/security/noInsecureRandom":                 {},
						"lint/security/noInsecureTlsConfig":              {},
						"lint/security/noLogInjectionTaint":              {},
						"lint/security/noPathTraversalTaint":             {},
						"lint/security/noPermissiveDirectoryPermissions": {},
						"lint/security/noPermissiveFilePermissions":      {},
						"lint/security/noPermissiveOsCreate":             {},
						"lint/security/noPermissiveWriteFilePermissions": {},
						"lint/security/noPredictableTempFile":            {},
						"lint/security/noSecretInSerialization":          {},
						"lint/security/noShortRsaKey":                    {},
						"lint/security/noSmtpInjectionTaint":             {},
						"lint/security/noSqlConcatenation":               {},
						"lint/security/noSqlFormatString":                {},
						"lint/security/noSqlInjectionTaint":              {},
						"lint/security/noSshAuthBypass":                  {},
						"lint/security/noSsrfTaint":                      {},
						"lint/security/noTemplateInjection":              {},
						"lint/security/noTlsSessionResumptionBypass":     {},
						"lint/security/noUnboundedFormParsing":           {},
						"lint/security/noUnsafeCorsBypass":               {},
						"lint/security/noUnsafeRedirectPolicy":           {},
						"lint/security/noUnescapedHtmlTemplate":          {},
						"lint/security/noUnsafeDeserialization":          {},
						"lint/security/noVariableCommandExecution":       {},
						"lint/security/noWeakCryptoHash":                 {},
						"lint/security/noWeakEncryptionAlgorithm":        {},
						"lint/security/noXssTaint":                       {},
						"lint/security/noZipSlip":                        {},
						"lint/style/noExcessiveBlankIdentifiers":         {},
						"lint/style/noDotImport":                         {},
						"lint/style/noInitFunction":                      {},
						"lint/style/noLeadingBlankLine":                  {},
						"lint/style/noTrailingBlankLine":                 {},
						"lint/style/noRepeatedStrings":                   {},
						"lint/style/noUnnecessaryConversion":             {},
						"lint/style/noUnnecessaryDeref":                  {},
						"lint/style/noUnnecessaryBlankIdentifier":        {},
						"lint/style/noUnnecessaryLoopVarCopy":            {},
						"lint/style/usePrintfSuffix":                     {},
						"lint/correctness/noSingleArgAppend":             {},
						"lint/correctness/noAsmDeclMismatch":             {},
						"lint/correctness/noSelfAssignment":              {},
						"lint/correctness/noSelfReferencingFinalizer":    {},
						"lint/correctness/noAtomicAssignMisuse":          {},
						"lint/correctness/noMalformedBuildTag":           {},
						"lint/correctness/noMalformedDirective":          {},
						"lint/correctness/noMalformedStructTag":          {},
						"lint/correctness/noCgoPointerViolation":         {},
						"lint/style/noUnkeyedLiteral":                    {},
						"lint/correctness/noCopiedLock":                  {},
						"lint/correctness/noDeferTimeMisuse":             {},
						"lint/correctness/noIncorrectTimeFormat":         {},
						"lint/correctness/noHttpClientGet":               {},
						"lint/correctness/noHttpClientHead":              {},
						"lint/correctness/noHttpClientPost":              {},
						"lint/correctness/noHttpClientPostForm":          {},
						"lint/correctness/noHttpGet":                     {},
						"lint/correctness/noHttpHead":                    {},
						"lint/correctness/noHttpPost":                    {},
						"lint/correctness/noHttpPostForm":                {},
						"lint/correctness/noHttpNewRequest":              {},
						"lint/correctness/noHttpResponseMisuse":          {},
						"lint/correctness/noHttptestNewRequest":          {},
						"lint/correctness/noInvalidBinaryArg":             {},
						"lint/correctness/noInvalidErrorsAs":             {},
						"lint/correctness/noFramePointerClobber":         {},
						"lint/correctness/useJoinHostPort":               {},
						"lint/correctness/noImpossibleInterfaceAssert":   {},
						"lint/correctness/noLoopClosureCapture":          {},
						"lint/correctness/noLostCancel":                  {},
						"lint/correctness/noWaitGroupMisuse":             {},
						"lint/correctness/noWriterBufferModification":    {},
						"lint/correctness/noNilContext":                  {},
						"lint/correctness/noNilFuncComparison":           {},
						"lint/correctness/noNilMapAssignment":            {},
						"lint/correctness/noPrintfFormatMismatch":        {},
						"lint/correctness/noExcessiveShift":              {},
						"lint/correctness/noUnbufferedSignalChannel":     {},
						"lint/correctness/noUntrappableSignal":           {},
						"lint/correctness/noSlogKeyValueMismatch":        {},
						"lint/correctness/noOddSizeSliceArg":             {},
						"lint/correctness/noStdMethodSignatureMismatch":  {},
						"lint/correctness/noStdlibVersionMismatch":       {},
						"lint/correctness/noStringIntConversion":         {},
						"lint/correctness/noTestFatalInGoroutine":        {},
						"lint/correctness/noTestMainWithoutExit":         {},
						"lint/correctness/noMalformedTestFunction":       {},
						"lint/correctness/noNonCanonicalHeaderKey":       {},
						"lint/style/noRedundantCanonicalHeaderKey":       {},
						"lint/correctness/noNonPointerUnmarshal":         {},
						"lint/correctness/noUnmarshalableType":           {},
						"lint/correctness/noUnreachableCode":             {},
						"lint/correctness/noInvalidUnsafePointer":        {},
						"lint/correctness/noUnusedFunctionResult":        {},
						"lint/correctness/noUnusedWrite":                 {},
						"lint/correctness/noIneffectualAssignment":       {},
						"lint/correctness/noIneffectiveBreak":            {},
						"lint/style/noLineTooLong":                       {},
						"lint/style/noMagicNumberInArgument":             {},
						"lint/style/noMagicNumberInAssignment":           {},
						"lint/style/noMagicNumberInCase":                 {},
						"lint/style/noMagicNumberInCondition":            {},
						"lint/style/noMagicNumberInOperation":            {},
						"lint/style/noMagicNumberInReturn":               {},
						"lint/correctness/noMisspelledWords":             {},
						"lint/correctness/noNanComparison":               {},
						"lint/style/noNakedReturn":                       {},
						"lint/correctness/noTlsConnHandshake":            {},
						"lint/correctness/noTlsDial":                     {},
						"lint/correctness/noTlsDialWithDialer":           {},
						"lint/correctness/noNetDial":                     {},
						"lint/correctness/noNetDialTimeout":              {},
						"lint/correctness/noNetIpBytesEqual":             {},
						"lint/correctness/noNetListen":                   {},
						"lint/correctness/noNetListenPacket":             {},
						"lint/correctness/noNetLookupCname":              {},
						"lint/correctness/noNetLookupHost":               {},
						"lint/correctness/noNetLookupIp":                 {},
						"lint/correctness/noNetLookupMx":                 {},
						"lint/correctness/noNetLookupPort":               {},
						"lint/correctness/noNetLookupSrv":                {},
						"lint/correctness/noNetLookupNs":                 {},
						"lint/correctness/noNetLookupTxt":                {},
						"lint/correctness/noNetLookupAddr":               {},
						"lint/correctness/noSqlDbBegin":                  {},
						"lint/correctness/noSqlDbExec":                   {},
						"lint/correctness/noSqlDbPing":                   {},
						"lint/correctness/noSqlDbPrepare":                {},
						"lint/correctness/noSqlDbQuery":                  {},
						"lint/correctness/noSqlDbQueryRow":               {},
						"lint/correctness/noSqlTxExec":                   {},
						"lint/correctness/noSqlTxPrepare":                {},
						"lint/correctness/noSqlTxQuery":                  {},
						"lint/correctness/noSqlTxQueryRow":               {},
						"lint/correctness/noSqlTxStmt":                   {},
						"lint/correctness/noSqlStmtExec":                 {},
						"lint/correctness/noSqlStmtQuery":                {},
						"lint/correctness/noSqlStmtQueryRow":             {},
						"lint/correctness/noExecCommand":                 {},
						"lint/correctness/noInvalidExecCommandArg":       {},
						"lint/style/noNolintNonMachineReadable":          {},
						"lint/correctness/noNolintParseError":            {},
						"lint/style/noNolintWithoutSpecificLinter":       {},
						"lint/style/noNolintWithoutExplanation":          {},
						"lint/suspicious/noUnusedNolint":                 {},
						"lint/suspicious/noUnusedParameter":              {},
						"lint/suspicious/noConstantParameter":            {},
						"lint/suspicious/noConstantResult":               {},
						"lint/correctness/noUnusedFunction":              {},
						"lint/correctness/noUnusedType":                  {},
						"lint/correctness/noUnusedVariable":              {},
						"lint/correctness/noUnusedConstant":              {},
						"lint/correctness/noUnusedField":                 {},
						"lint/suspicious/noMismatchedAppendAssign":       {},
						"lint/performance/noHugeParam":                   {},
						"lint/performance/noStringIndexAllocation":       {},
						"lint/performance/useStringMapKey":               {},
						"lint/performance/noRedundantStringByteConversion": {},
						"lint/performance/useInlineMathPow":             {},
						"lint/performance/useCombinedAppend":             {},
						"lint/suspicious/noSwappedArguments":             {},
						"lint/style/useAssignmentOperator":               {},
						"lint/correctness/noNoopFunctionCall":            {},
						"lint/correctness/noImpossibleCondition":         {},
						"lint/style/noCapitalizedLocal":                  {},
						"lint/style/noCapitalizedErrorString":            {},
						"lint/correctness/noUnreachableTypeCase":         {},
						"lint/style/noMisplacedDefaultCase":              {},
						"lint/style/useStandardCodegenComment":           {},
						"lint/style/useStandardDeprecationComment":       {},
						"lint/style/useCommentSpacing":                   {},
						"lint/correctness/noDuplicateBuildConstraint":    {},
						"lint/correctness/noDuplicateArgument":           {},
						"lint/suspicious/noDuplicateBranchBody":          {},
						"lint/correctness/noDuplicateCase":               {},
						"lint/suspicious/noDuplicateSubExpression":       {},
						"lint/suspicious/noDuplicateCutsetChars":         {},
						"lint/suspicious/noDuplicateIfCondition":         {},
						"lint/style/useElseIf":                           {},
						"lint/correctness/noExitAfterDefer":              {},
						"lint/correctness/noFlagDerefBeforeParse":        {},
						"lint/correctness/noInvalidFlagName":             {},
						"lint/correctness/noInvalidHostPort":             {},
						"lint/style/useSwitch":                           {},
						"lint/style/useTaggedSwitch":                     {},
						"lint/style/useSimplifiedBoolReturn":             {},
						"lint/style/useSimplifiedPrintFormat":            {},
						"lint/style/useStringConversionInPrint":          {},
						"lint/style/usePackageComment":                   {},
						"lint/style/useIdiomaticErrorName":               {},
						"lint/style/useIdiomaticNaming":                  {},
						"lint/style/useIdiomaticReceiverName":            {},
						"lint/style/useIdiomaticDurationName":            {},
						"lint/style/useConsistentReceiverName":           {},
						"lint/style/noControlCharInString":               {},
						"lint/style/useFuncDocPrefix":                    {},
						"lint/style/useTypeDocPrefix":                    {},
						"lint/style/useVarConstDocPrefix":                {},
						"lint/style/noRedundantMakeArgs":                 {},
						"lint/style/noRedundantNilLoopCheck":             {},
						"lint/style/noRedundantNilSliceCheck":            {},
						"lint/style/noRedundantNilTypeCheck":             {},
						"lint/style/noDefaultSliceIndex":                 {},
						"lint/style/useSliceAppend":                      {},
						"lint/correctness/useTimeEqual":                  {},
						"lint/style/useTimeSince":                        {},
						"lint/style/useTimeSleep":                        {},
						"lint/style/useTimeUntil":                        {},
						"lint/style/useTypeConversion":                   {},
						"lint/style/useTrimFunction":                     {},
						"lint/style/useMergedConditionalDecl":            {},
						"lint/style/useMergedVarDecl":                    {},
						"lint/style/noRedundantControlFlow":              {},
						"lint/style/useErrorMethod":                      {},
						"lint/style/useErrorLastReturn":                  {},
						"lint/style/useDirectStringRange":                {},
						"lint/style/useBufferStringOrBytes":              {},
						"lint/style/noGuardAroundDelete":                 {},
						"lint/style/noGuardAroundMapAccess":              {},
						"lint/suspicious/noSuspiciousMapKey":             {},
						"lint/suspicious/noSuspiciousTimeSleep":          {},
						"lint/suspicious/noEmptyBranch":                  {},
						"lint/suspicious/noEmptyCriticalSection":         {},
						"lint/correctness/noTimeTick":                    {},
						"lint/correctness/noTimerResetRetval":            {},
						"lint/style/noImmediateNewDeref":                 {},
						"lint/correctness/noOffByOneError":               {},
						"lint/performance/noRangeExprCopy":               {},
						"lint/performance/noRangeValCopy":                {},
						"lint/performance/noRedundantRuneConversion":     {},
						"lint/style/useRegexpMustCompile":                {},
						"lint/style/useRawStringRegexp":                  {},
						"lint/performance/useCompiledRegexp":             {},
						"lint/style/noSingleCaseSelect":                  {},
						"lint/style/noSingleCaseSwitch":                  {},
						"lint/style/noRedundantSwitchTrue":               {},
						"lint/suspicious/noAlwaysTrueLenCheck":           {},
						"lint/style/noRedundantTypeAssertion":            {},
						"lint/style/useTypeSwitchGuard":                  {},
						"lint/style/useTypeAssertResult":                 {},
						"lint/style/noUnnecessaryLambda":                 {},
						"lint/style/noRedundantSliceExpression":          {},
						"lint/style/useParallelAssignSwap":               {},
						"lint/style/useConvenienceFunc":                  {},
						"lint/suspicious/noZeroBytesRepeat":              {},
						"lint/correctness/noInvalidRegexp":               {},
						"lint/correctness/noInvalidTemplate":             {},
						"lint/correctness/noInvalidUrlParse":             {},
						"lint/correctness/noInvalidUtf8StringArg":        {},
						"lint/correctness/noInappropriateContextKey":     {},
						"lint/correctness/noInvalidStrconvArg":           {},
						"lint/correctness/noOverlappingEncoderSlice":     {},
						"lint/correctness/noBenchmarkNAssignment":        {},
						"lint/suspicious/noAddressOfDereference":         {},
						"lint/suspicious/noAddressNilComparison":         {},
						"lint/suspicious/noUnsignedNegativeComparison":   {},
						"lint/suspicious/noUnobservedFieldAssign":        {},
						"lint/suspicious/noOverwrittenArgument":          {},
						"lint/suspicious/noSingleIterationLoop":          {},
						"lint/suspicious/noInvariantLoopCondition":       {},
						"lint/suspicious/noUselessMathCall":              {},
						"lint/suspicious/noIneffectiveBitwiseOp":         {},
						"lint/suspicious/noIneffectiveRandCall":          {},
						"lint/suspicious/noImpossibleNilComparison":      {},
						"lint/suspicious/noImpossibleBuiltinResult":      {},
						"lint/suspicious/noIntegerDivisionTruncation":    {},
						"lint/suspicious/noImpreciseConstant":            {},
						"lint/suspicious/noImplicitConstValue":           {},
						"lint/suspicious/noIgnoredQueryModification":     {},
						"lint/suspicious/noModuloOne":                    {},
						"lint/suspicious/noNeverNilCheck":                {},
						"lint/suspicious/noNonOctalFileMode":             {},
						"lint/suspicious/noUnmarshalableStruct":          {},
						"lint/suspicious/noDubiousBitShift":              {},
						"lint/suspicious/noTempDirDeletion":              {},
						"lint/suspicious/noTypeAssertElseMisread":        {},
						"lint/style/useCopyBuiltin":                      {},
						"lint/style/useCopyForSlide":                     {},
						"lint/style/noExplicitBoolComparison":            {},
						"lint/style/useStringsContains":                  {},
						"lint/style/useBytesEqual":                       {},
						"lint/style/useInfiniteFor":                      {},
						"lint/style/noRedundantVarType":                  {},
						"lint/complexity/useLoopCondition":                {},
						"lint/style/useSimplifiedSelector":                {},
						"lint/style/usePackageComments":                               {},
						"lint/style/noRedundantRangeVal":                  {},
						"lint/style/useReceiverNaming":                                {},
						"lint/style/noSuperfluousElse":                               {},
						"lint/style/useTimeNaming":                       {},
						"lint/style/noUnexportedReturn":                  {},
						"lint/style/noVarDeclaration":                    {},
						"lint/style/useVarNaming":                                     {},
						"lint/performance/usePointerInSyncPool":          {},
					},
				},
			},
			"config with non-defaults": {
				confPath: "non-defaults.toml",
				wantConfig: lint.Config{
					Confidence:            0.5,
					Severity:              lint.SeverityError,
					IgnoreGeneratedHeader: true,
					EnableDefaultRules:    true,
					ErrorCode:             2,
					WarningCode:           1,
					Rules: lint.RulesConfig{
						"lint/complexity/noExcessiveArguments": {
							Severity: lint.SeverityWarning,
							Exclude:  []string{"excluded/file.go"},
							Arguments: lint.Arguments{
								[]any{4},
							},
						},
						"lint/style/noBlankImport": {
							Disabled: true,
							Severity: lint.SeverityError,
						},
						"lint/style/useContextAsFirstParam": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noContextKeysType": {
							Severity: lint.SeverityError,
						},
						"empty-block": {
							Severity: lint.SeverityError,
						},
						"lint/style/useErrorNaming": {
							Severity: lint.SeverityError,
						},
						"lint/style/noErrorStrings": {
							Severity: lint.SeverityError,
						},
						"lint/style/useErrorf": {
							Severity: lint.SeverityError,
						},
						"lint/style/useExportedComment": {
							Severity: lint.SeverityError,
							Arguments: lint.Arguments{
								"check-private-receivers", "disable-stuttering-check",
							},
							Exclude: []string{"excluded/file-exported.go"},
						},
						"increment-decrement": {
							Severity: lint.SeverityError,
						},
						"lint/style/useIndentErrorFlow": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnclosedBodies": {
							Severity: lint.SeverityError,
						},
						"lint/complexity/noDuplicateCode": {
							Severity: lint.SeverityError,
						},
						"lint/complexity/noDuplicateCodeV2": {
							Severity: lint.SeverityError,
						},
						"lint/complexity/noExcessiveStatements": {
							Severity: lint.SeverityError,
						},
						"lint/complexity/noHighCyclomaticComplexity": {
							Severity: lint.SeverityError,
						},
						"lint/complexity/noLongFunctions": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noDeniedImport": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noDeprecatedUsage": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noDirectErrorComparison": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noDiscardedAppend": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noDeferCloseBeforeErrCheck": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noDeferInInfiniteLoop": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noDynamicErrors": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noEmptyForLoop": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSelectBreakConfusion": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noInfiniteRecursion": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noFileScopedDeniedImport": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSpaceInDirective": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnallowedImport": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUncheckedError": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noRangeVariableAlias": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSliceBoundsOutOfRange": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnrecognizedDirective": {
							Severity: lint.SeverityError,
						},
						"lint/security/noAtoiOverflow": {
							Severity: lint.SeverityError,
						},
						"lint/security/noBindToAllInterfaces": {
							Severity: lint.SeverityError,
						},
						"lint/security/noExposedPprof": {
							Severity: lint.SeverityError,
						},
						"lint/security/noFilesystemRootServing": {
							Severity: lint.SeverityError,
						},
						"lint/security/noHardcodedCredentials": {
							Severity: lint.SeverityError,
						},
						"lint/security/noHardcodedIv": {
							Severity: lint.SeverityError,
						},
						"lint/security/noHttpRequestSmuggling": {
							Severity: lint.SeverityError,
						},
						"lint/security/noInsecureCookie": {
							Severity: lint.SeverityError,
						},
						"lint/security/noInsecureHostKeyCallback": {
							Severity: lint.SeverityError,
						},
						"lint/security/noIntegerOverflowConversion": {
							Severity: lint.SeverityError,
						},
						"lint/security/noMissingReadHeaderTimeout": {
							Severity: lint.SeverityError,
						},
						"lint/security/noServeWithoutTimeout": {
							Severity: lint.SeverityError,
						},
						"lint/security/noSsrfViaVariable": {
							Severity: lint.SeverityError,
						},
						"lint/security/noTrojanSourceBidi": {
							Severity: lint.SeverityError,
						},
						"lint/security/noUnboundedDecompression": {
							Severity: lint.SeverityError,
						},
						"lint/security/noUnsafePackage": {
							Severity: lint.SeverityError,
						},
						"lint/security/noCgiImport": {
							Severity: lint.SeverityError,
						},
						"lint/security/noCommandInjectionTaint": {
							Severity: lint.SeverityError,
						},
						"lint/security/noContextPropagationFailure": {
							Severity: lint.SeverityError,
						},
						"lint/security/noDeprecatedHashFunction": {
							Severity: lint.SeverityError,
						},
						"lint/security/noFileInclusionViaVariable": {
							Severity: lint.SeverityError,
						},
						"lint/security/noFilesystemToctou": {
							Severity: lint.SeverityError,
						},
						"lint/security/noInsecureRandom": {
							Severity: lint.SeverityError,
						},
						"lint/security/noInsecureTlsConfig": {
							Severity: lint.SeverityError,
						},
						"lint/security/noLogInjectionTaint": {
							Severity: lint.SeverityError,
						},
						"lint/security/noPathTraversalTaint": {
							Severity: lint.SeverityError,
						},
						"lint/security/noPermissiveDirectoryPermissions": {
							Severity: lint.SeverityError,
						},
						"lint/security/noPermissiveFilePermissions": {
							Severity: lint.SeverityError,
						},
						"lint/security/noPermissiveOsCreate": {
							Severity: lint.SeverityError,
						},
						"lint/security/noPermissiveWriteFilePermissions": {
							Severity: lint.SeverityError,
						},
						"lint/security/noPredictableTempFile": {
							Severity: lint.SeverityError,
						},
						"lint/security/noSecretInSerialization": {
							Severity: lint.SeverityError,
						},
						"lint/security/noShortRsaKey": {
							Severity: lint.SeverityError,
						},
						"lint/security/noSmtpInjectionTaint": {
							Severity: lint.SeverityError,
						},
						"lint/security/noSqlConcatenation": {
							Severity: lint.SeverityError,
						},
						"lint/security/noSqlFormatString": {
							Severity: lint.SeverityError,
						},
						"lint/security/noSqlInjectionTaint": {
							Severity: lint.SeverityError,
						},
						"lint/security/noSshAuthBypass": {
							Severity: lint.SeverityError,
						},
						"lint/security/noSsrfTaint": {
							Severity: lint.SeverityError,
						},
						"lint/security/noTemplateInjection": {
							Severity: lint.SeverityError,
						},
						"lint/security/noTlsSessionResumptionBypass": {
							Severity: lint.SeverityError,
						},
						"lint/security/noUnboundedFormParsing": {
							Severity: lint.SeverityError,
						},
						"lint/security/noUnsafeCorsBypass": {
							Severity: lint.SeverityError,
						},
						"lint/security/noUnsafeRedirectPolicy": {
							Severity: lint.SeverityError,
						},
						"lint/security/noUnescapedHtmlTemplate": {
							Severity: lint.SeverityError,
						},
						"lint/security/noUnsafeDeserialization": {
							Severity: lint.SeverityError,
						},
						"lint/security/noVariableCommandExecution": {
							Severity: lint.SeverityError,
						},
						"lint/security/noWeakCryptoHash": {
							Severity: lint.SeverityError,
						},
						"lint/security/noWeakEncryptionAlgorithm": {
							Severity: lint.SeverityError,
						},
						"lint/security/noXssTaint": {
							Severity: lint.SeverityError,
						},
						"lint/security/noZipSlip": {
							Severity: lint.SeverityError,
						},
						"lint/style/noExcessiveBlankIdentifiers": {
							Severity: lint.SeverityError,
						},
						"lint/style/noDotImport": {
							Severity: lint.SeverityError,
						},
						"lint/style/noInitFunction": {
							Severity: lint.SeverityError,
						},
						"lint/style/noLeadingBlankLine": {
							Severity: lint.SeverityError,
						},
						"lint/style/noTrailingBlankLine": {
							Severity: lint.SeverityError,
						},
						"lint/style/noRepeatedStrings": {
							Severity: lint.SeverityError,
						},
						"lint/style/noUnnecessaryConversion": {
							Severity: lint.SeverityError,
						},
						"lint/style/noUnnecessaryDeref": {
							Severity: lint.SeverityError,
						},
						"lint/style/noUnnecessaryBlankIdentifier": {
							Severity: lint.SeverityError,
						},
						"lint/style/noUnnecessaryLoopVarCopy": {
							Severity: lint.SeverityError,
						},
						"lint/style/usePrintfSuffix": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSingleArgAppend": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noAsmDeclMismatch": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSelfAssignment": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSelfReferencingFinalizer": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noAtomicAssignMisuse": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noMalformedBuildTag": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noMalformedDirective": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noMalformedStructTag": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noCgoPointerViolation": {
							Severity: lint.SeverityError,
						},
						"lint/style/noUnkeyedLiteral": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noCopiedLock": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noDeferTimeMisuse": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noIncorrectTimeFormat": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noHttpClientGet": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noHttpClientHead": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noHttpClientPost": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noHttpClientPostForm": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noHttpGet": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noHttpHead": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noHttpPost": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noHttpPostForm": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noHttpNewRequest": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noHttpResponseMisuse": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noHttptestNewRequest": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noInvalidBinaryArg": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noInvalidErrorsAs": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noFramePointerClobber": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/useJoinHostPort": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noImpossibleInterfaceAssert": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noLoopClosureCapture": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noLostCancel": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noWaitGroupMisuse": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noWriterBufferModification": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNilContext": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNilFuncComparison": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNilMapAssignment": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noPrintfFormatMismatch": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noExcessiveShift": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnbufferedSignalChannel": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUntrappableSignal": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSlogKeyValueMismatch": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noOddSizeSliceArg": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noStdMethodSignatureMismatch": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noStdlibVersionMismatch": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noStringIntConversion": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noTestFatalInGoroutine": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noTestMainWithoutExit": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noMalformedTestFunction": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNonCanonicalHeaderKey": {
							Severity: lint.SeverityError,
						},
						"lint/style/noRedundantCanonicalHeaderKey": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNonPointerUnmarshal": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnmarshalableType": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnreachableCode": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noInvalidUnsafePointer": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnusedFunctionResult": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnusedWrite": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noIneffectualAssignment": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noIneffectiveBreak": {
							Severity: lint.SeverityError,
						},
						"lint/style/noLineTooLong": {
							Severity: lint.SeverityError,
						},
						"lint/style/noMagicNumberInArgument": {
							Severity: lint.SeverityError,
						},
						"lint/style/noMagicNumberInAssignment": {
							Severity: lint.SeverityError,
						},
						"lint/style/noMagicNumberInCase": {
							Severity: lint.SeverityError,
						},
						"lint/style/noMagicNumberInCondition": {
							Severity: lint.SeverityError,
						},
						"lint/style/noMagicNumberInOperation": {
							Severity: lint.SeverityError,
						},
						"lint/style/noMagicNumberInReturn": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noMisspelledWords": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNanComparison": {
							Severity: lint.SeverityError,
						},
						"lint/style/noNakedReturn": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noTlsConnHandshake": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noTlsDial": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noTlsDialWithDialer": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNetDial": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNetDialTimeout": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNetIpBytesEqual": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNetListen": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNetListenPacket": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNetLookupCname": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNetLookupHost": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNetLookupIp": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNetLookupMx": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNetLookupPort": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNetLookupSrv": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNetLookupNs": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNetLookupTxt": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNetLookupAddr": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSqlDbBegin": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSqlDbExec": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSqlDbPing": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSqlDbPrepare": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSqlDbQuery": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSqlDbQueryRow": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSqlTxExec": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSqlTxPrepare": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSqlTxQuery": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSqlTxQueryRow": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSqlTxStmt": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSqlStmtExec": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSqlStmtQuery": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noSqlStmtQueryRow": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noExecCommand": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noInvalidExecCommandArg": {
							Severity: lint.SeverityError,
						},
						"lint/style/noNolintNonMachineReadable": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNolintParseError": {
							Severity: lint.SeverityError,
						},
						"lint/style/noNolintWithoutSpecificLinter": {
							Severity: lint.SeverityError,
						},
						"lint/style/noNolintWithoutExplanation": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noUnusedNolint": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noUnusedParameter": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noConstantParameter": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noConstantResult": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnusedFunction": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnusedType": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnusedVariable": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnusedConstant": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnusedField": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noMismatchedAppendAssign": {
							Severity: lint.SeverityError,
						},
						"lint/performance/noHugeParam": {
							Severity: lint.SeverityError,
						},
						"lint/performance/noStringIndexAllocation": {
							Severity: lint.SeverityError,
						},
						"lint/performance/useStringMapKey": {
							Severity: lint.SeverityError,
						},
						"lint/performance/noRedundantStringByteConversion": {
							Severity: lint.SeverityError,
						},
						"lint/performance/useInlineMathPow": {
							Severity: lint.SeverityError,
						},
						"lint/performance/useCombinedAppend": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noSwappedArguments": {
							Severity: lint.SeverityError,
						},
						"lint/style/useAssignmentOperator": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNoopFunctionCall": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noImpossibleCondition": {
							Severity: lint.SeverityError,
						},
						"lint/style/noCapitalizedLocal": {
							Severity: lint.SeverityError,
						},
						"lint/style/noCapitalizedErrorString": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnreachableTypeCase": {
							Severity: lint.SeverityError,
						},
						"lint/style/noMisplacedDefaultCase": {
							Severity: lint.SeverityError,
						},
						"lint/style/useStandardCodegenComment": {
							Severity: lint.SeverityError,
						},
						"lint/style/useStandardDeprecationComment": {
							Severity: lint.SeverityError,
						},
						"lint/style/useCommentSpacing": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noDuplicateBuildConstraint": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noDuplicateArgument": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noDuplicateBranchBody": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noDuplicateCase": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noDuplicateSubExpression": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noDuplicateCutsetChars": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noDuplicateIfCondition": {
							Severity: lint.SeverityError,
						},
						"lint/style/useElseIf": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noExitAfterDefer": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noFlagDerefBeforeParse": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noInvalidFlagName": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noInvalidHostPort": {
							Severity: lint.SeverityError,
						},
						"lint/style/useSwitch": {
							Severity: lint.SeverityError,
						},
						"lint/style/useTaggedSwitch": {
							Severity: lint.SeverityError,
						},
						"lint/style/useSimplifiedBoolReturn": {
							Severity: lint.SeverityError,
						},
						"lint/style/useSimplifiedPrintFormat": {
							Severity: lint.SeverityError,
						},
						"lint/style/useStringConversionInPrint": {
							Severity: lint.SeverityError,
						},
						"lint/style/usePackageComment": {
							Severity: lint.SeverityError,
						},
						"lint/style/useIdiomaticErrorName": {
							Severity: lint.SeverityError,
						},
						"lint/style/useIdiomaticNaming": {
							Severity: lint.SeverityError,
						},
						"lint/style/useIdiomaticReceiverName": {
							Severity: lint.SeverityError,
						},
						"lint/style/useIdiomaticDurationName": {
							Severity: lint.SeverityError,
						},
						"lint/style/useConsistentReceiverName": {
							Severity: lint.SeverityError,
						},
						"lint/style/noControlCharInString": {
							Severity: lint.SeverityError,
						},
						"lint/style/useFuncDocPrefix": {
							Severity: lint.SeverityError,
						},
						"lint/style/useTypeDocPrefix": {
							Severity: lint.SeverityError,
						},
						"lint/style/useVarConstDocPrefix": {
							Severity: lint.SeverityError,
						},
						"lint/style/noRedundantMakeArgs": {
							Severity: lint.SeverityError,
						},
						"lint/style/noRedundantNilLoopCheck": {
							Severity: lint.SeverityError,
						},
						"lint/style/noRedundantNilSliceCheck": {
							Severity: lint.SeverityError,
						},
						"lint/style/noRedundantNilTypeCheck": {
							Severity: lint.SeverityError,
						},
						"lint/style/noDefaultSliceIndex": {
							Severity: lint.SeverityError,
						},
						"lint/style/useSliceAppend": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/useTimeEqual": {
							Severity: lint.SeverityError,
						},
						"lint/style/useTimeSince": {
							Severity: lint.SeverityError,
						},
						"lint/style/useTimeSleep": {
							Severity: lint.SeverityError,
						},
						"lint/style/useTimeUntil": {
							Severity: lint.SeverityError,
						},
						"lint/style/useTypeConversion": {
							Severity: lint.SeverityError,
						},
						"lint/style/useTrimFunction": {
							Severity: lint.SeverityError,
						},
						"lint/style/useMergedConditionalDecl": {
							Severity: lint.SeverityError,
						},
						"lint/style/useMergedVarDecl": {
							Severity: lint.SeverityError,
						},
						"lint/style/noRedundantControlFlow": {
							Severity: lint.SeverityError,
						},
						"lint/style/useErrorMethod": {
							Severity: lint.SeverityError,
						},
						"lint/style/useErrorLastReturn": {
							Severity: lint.SeverityError,
						},
						"lint/style/useDirectStringRange": {
							Severity: lint.SeverityError,
						},
						"lint/style/useBufferStringOrBytes": {
							Severity: lint.SeverityError,
						},
						"lint/style/noGuardAroundDelete": {
							Severity: lint.SeverityError,
						},
						"lint/style/noGuardAroundMapAccess": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noSuspiciousMapKey": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noSuspiciousTimeSleep": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noEmptyBranch": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noEmptyCriticalSection": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noTimeTick": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noTimerResetRetval": {
							Severity: lint.SeverityError,
						},
						"lint/style/noImmediateNewDeref": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noOffByOneError": {
							Severity: lint.SeverityError,
						},
						"lint/performance/noRangeExprCopy": {
							Severity: lint.SeverityError,
						},
						"lint/performance/noRangeValCopy": {
							Severity: lint.SeverityError,
						},
						"lint/performance/noRedundantRuneConversion": {
							Severity: lint.SeverityError,
						},
						"lint/style/useRegexpMustCompile": {
							Severity: lint.SeverityError,
						},
						"lint/style/useRawStringRegexp": {
							Severity: lint.SeverityError,
						},
						"lint/performance/useCompiledRegexp": {
							Severity: lint.SeverityError,
						},
						"lint/style/noSingleCaseSelect": {
							Severity: lint.SeverityError,
						},
						"lint/style/noSingleCaseSwitch": {
							Severity: lint.SeverityError,
						},
						"lint/style/noRedundantSwitchTrue": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noAlwaysTrueLenCheck": {
							Severity: lint.SeverityError,
						},
						"lint/style/noRedundantTypeAssertion": {
							Severity: lint.SeverityError,
						},
						"lint/style/useTypeSwitchGuard": {
							Severity: lint.SeverityError,
						},
						"lint/style/useTypeAssertResult": {
							Severity: lint.SeverityError,
						},
						"lint/style/noUnnecessaryLambda": {
							Severity: lint.SeverityError,
						},
						"lint/style/noRedundantSliceExpression": {
							Severity: lint.SeverityError,
						},
						"lint/style/useParallelAssignSwap": {
							Severity: lint.SeverityError,
						},
						"lint/style/useConvenienceFunc": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noZeroBytesRepeat": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noInvalidRegexp": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noInvalidTemplate": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noInvalidUrlParse": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noInvalidUtf8StringArg": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noInappropriateContextKey": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noInvalidStrconvArg": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noOverlappingEncoderSlice": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noBenchmarkNAssignment": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noAddressOfDereference": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noAddressNilComparison": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noUnsignedNegativeComparison": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noUnobservedFieldAssign": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noOverwrittenArgument": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noSingleIterationLoop": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noInvariantLoopCondition": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noUselessMathCall": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noIneffectiveBitwiseOp": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noIneffectiveRandCall": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noImpossibleNilComparison": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noImpossibleBuiltinResult": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noIntegerDivisionTruncation": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noImpreciseConstant": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noImplicitConstValue": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noIgnoredQueryModification": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noModuloOne": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noNeverNilCheck": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noNonOctalFileMode": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noUnmarshalableStruct": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noDubiousBitShift": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noTempDirDeletion": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noTypeAssertElseMisread": {
							Severity: lint.SeverityError,
						},
						"lint/style/useCopyBuiltin": {
							Severity: lint.SeverityError,
						},
						"lint/style/useCopyForSlide": {
							Severity: lint.SeverityError,
						},
						"lint/style/noExplicitBoolComparison": {
							Severity: lint.SeverityError,
						},
						"lint/style/useStringsContains": {
							Severity: lint.SeverityError,
						},
						"lint/style/useBytesEqual": {
							Severity: lint.SeverityError,
						},
						"lint/style/useInfiniteFor": {
							Severity: lint.SeverityError,
						},
						"lint/style/noRedundantVarType": {
							Severity: lint.SeverityError,
						},
						"lint/complexity/useLoopCondition": {
							Severity: lint.SeverityError,
						},
						"lint/style/useSimplifiedSelector": {
							Severity: lint.SeverityError,
						},
						"lint/style/usePackageComments": {
							Severity: lint.SeverityError,
						},
						"lint/style/noRedundantRangeVal": {
							Severity: lint.SeverityError,
						},
						"lint/style/useReceiverNaming": {
							Severity: lint.SeverityError,
						},
						"lint/style/noSuperfluousElse": {
							Severity: lint.SeverityError,
						},
						"lint/style/useTimeNaming": {
							Severity: lint.SeverityError,
						},
						"lint/style/noUnexportedReturn": {
							Severity: lint.SeverityError,
						},
						"lint/style/noVarDeclaration": {
							Severity: lint.SeverityError,
						},
						"lint/style/useVarNaming": {
							Severity: lint.SeverityError,
						},
						"lint/performance/usePointerInSyncPool": {
							Severity: lint.SeverityError,
						},
					},
				},
			},
		} {
			t.Run(name, func(t *testing.T) {
				var cfgPath string
				if tc.confPath != "" {
					cfgPath = filepath.Join("testdata", tc.confPath)
				}

				cfg, err := config.GetConfig(cfgPath)
				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}
				if cfg.IgnoreGeneratedHeader != tc.wantConfig.IgnoreGeneratedHeader {
					t.Errorf("IgnoreGeneratedHeader: expected %v, got %v", tc.wantConfig.IgnoreGeneratedHeader, cfg.IgnoreGeneratedHeader)
				}
				if cfg.Confidence != tc.wantConfig.Confidence {
					t.Errorf("Confidence: expected %v, got %v", tc.wantConfig.Confidence, cfg.Confidence)
				}
				if cfg.Severity != tc.wantConfig.Severity {
					t.Errorf("Severity: expected %v, got %v", tc.wantConfig.Severity, cfg.Severity)
				}
				if cfg.EnableAllRules != tc.wantConfig.EnableAllRules {
					t.Errorf("EnableAllRules: expected %v, got %v", tc.wantConfig.EnableAllRules, cfg.EnableAllRules)
				}
				if cfg.EnableDefaultRules != tc.wantConfig.EnableDefaultRules {
					t.Errorf("EnableDefaultRules: expected %v, got %v", tc.wantConfig.EnableDefaultRules, cfg.EnableDefaultRules)
				}
				if cfg.ErrorCode != tc.wantConfig.ErrorCode {
					t.Errorf("ErrorCode: expected %v, got %v", tc.wantConfig.ErrorCode, cfg.ErrorCode)
				}
				if cfg.WarningCode != tc.wantConfig.WarningCode {
					t.Errorf("WarningCode: expected %v, got %v", tc.wantConfig.WarningCode, cfg.WarningCode)
				}
				if !tc.wantConfig.GoVersion.Equal(cfg.GoVersion) {
					t.Errorf("GoVersion: expected %v, got %v", tc.wantConfig.GoVersion, cfg.GoVersion)
				}

				if len(cfg.Exclude) != len(tc.wantConfig.Exclude) {
					t.Errorf("Exclude length: expected %v, got %v", len(tc.wantConfig.Exclude), len(cfg.Exclude))
				} else {
					for i, exclude := range tc.wantConfig.Exclude {
						if cfg.Exclude[i] != exclude {
							t.Errorf("Exclude[%d]: expected %v, got %v", i, exclude, cfg.Exclude[i])
						}
					}
				}

				if len(cfg.Rules) != len(tc.wantConfig.Rules) {
					t.Errorf("Rules count: expected %v, got %v", len(tc.wantConfig.Rules), len(cfg.Rules))
				}
				for ruleName, wantRule := range tc.wantConfig.Rules {
					gotRule, exists := cfg.Rules[ruleName]
					if !exists {
						t.Errorf("Rule %q: expected to exist, but not found", ruleName)
						continue
					}
					if gotRule.Disabled != wantRule.Disabled {
						t.Errorf("Rule %q Disabled: expected %v, got %v", ruleName, wantRule.Disabled, gotRule.Disabled)
					}
					if gotRule.Severity != wantRule.Severity {
						t.Errorf("Rule %q Severity: expected %v, got %v", ruleName, wantRule.Severity, gotRule.Severity)
					}
					if len(gotRule.Arguments) != len(wantRule.Arguments) {
						t.Errorf("Rule %q Arguments length: expected %v, got %v", ruleName, len(wantRule.Arguments), len(gotRule.Arguments))
					}
					if len(gotRule.Exclude) != len(wantRule.Exclude) {
						t.Errorf("Rule %q Exclude length: expected %v, got %v", ruleName, len(wantRule.Exclude), len(gotRule.Exclude))
					} else {
						for i, wantExclude := range wantRule.Exclude {
							if gotRule.Exclude[i] != wantExclude {
								t.Errorf("Rule %q Exclude[%d]: expected %v, got %v", ruleName, i, wantExclude, gotRule.Exclude[i])
							}
						}
					}
				}
				// Check for unexpected rules in actual config
				for ruleName := range cfg.Rules {
					if _, exists := tc.wantConfig.Rules[ruleName]; !exists {
						t.Errorf("Rule %q: found in actual config but not expected", ruleName)
					}
				}

				if len(cfg.Directives) != len(tc.wantConfig.Directives) {
					t.Errorf("Directives count: expected %v, got %v", len(tc.wantConfig.Directives), len(cfg.Directives))
				}
				for directiveName, wantDirective := range tc.wantConfig.Directives {
					gotDirective, exists := cfg.Directives[directiveName]
					if !exists {
						t.Errorf("Directive %q: expected to exist, but not found", directiveName)
						continue
					}
					if gotDirective.Severity != wantDirective.Severity {
						t.Errorf("Directive %q Severity: expected %v, got %v", directiveName, wantDirective.Severity, gotDirective.Severity)
					}
				}
				// Check for unexpected directives in actual config
				for directiveName := range cfg.Directives {
					if _, exists := tc.wantConfig.Directives[directiveName]; !exists {
						t.Errorf("Directive %q: found in actual config but not expected", directiveName)
					}
				}
			})
		}

		t.Run("rule-level file filter excludes", func(t *testing.T) {
			cfg, err := config.GetConfig("testdata/rule-level-exclude-850.toml")
			if err != nil {
				t.Fatal("should be valid config")
			}
			r1 := cfg.Rules["r1"]
			if len(r1.Exclude) > 0 {
				t.Fatal("r1 should have empty excludes")
			}
			r2 := cfg.Rules["r2"]
			if len(r2.Exclude) != 1 {
				t.Fatal("r2 should have exclude set")
			}
			if !r2.MustExclude("some/file.go") {
				t.Fatal("r2 should be initialized and exclude some/file.go")
			}
			if r2.MustExclude("some/any-other.go") {
				t.Fatal("r2 should not exclude some/any-other.go")
			}
		})
	})

	t.Run("failure", func(t *testing.T) {
		for name, tc := range map[string]struct {
			confPath  string
			wantError string
		}{
			"unknown file": {
				confPath:  "unknown",
				wantError: "cannot read the config file",
			},
			"malformed file": {
				confPath:  "malformed.toml",
				wantError: "cannot parse the config file",
			},
			"invalid exclude pattern": {
				confPath:  "invalidExcludePattern.toml",
				wantError: "error in config of rule [lint/style/useVarNaming]",
			},
			"enableAllRules and enableDefaultRules both set": {
				confPath:  "enableAllAndDefault.toml",
				wantError: "config options enableAllRules and enableDefaultRules cannot be combined",
			},
		} {
			t.Run(name, func(t *testing.T) {
				_, err := config.GetConfig(filepath.Join("testdata", tc.confPath))

				if err != nil && !strings.Contains(err.Error(), tc.wantError) {
					t.Errorf("Unexpected error: want %q, got: %q", tc.wantError, err)
				}
			})
		}
	})
}

func TestGetLintingRules(t *testing.T) {
	const (
		// len of defaultRules
		defaultRulesCount = 350
		// len of allRules: update this when adding new rules
		allRulesCount = 493
	)

	tt := map[string]struct {
		confPath          string
		wantRulesCount    int
		wantEnabledRules  []string
		wantDisabledRules []string
		wantErr           string
	}{
		"no rules": {
			confPath:       "noRules.toml",
			wantRulesCount: 0,
			wantDisabledRules: []string{
				"noVarDeclaration",  // default rule
				"usePackageComments", // default rule
				"noDeepExit",        // non-default rule
			},
		},
		"enableAllRules without disabled rules": {
			confPath:       "enableAll.toml",
			wantRulesCount: allRulesCount,
			wantEnabledRules: []string{
				"noVarDeclaration",  // default rule
				"usePackageComments", // default rule
				"noDeepExit",        // non-default rule
			},
		},
		"enableAllRules with 2 disabled rules": {
			confPath:       "enableAllBut2.toml",
			wantRulesCount: allRulesCount - 2,
			wantEnabledRules: []string{
				"noVarDeclaration",  // default rule
				"usePackageComments", // default rule
				"noDeepExit",        // non-default rule
			},
			wantDisabledRules: []string{
				"useExportedComment",   // default rule
				"noExcessiveArguments", // non-default rule
			},
		},
		"enableDefaultRules without disabled rules": {
			confPath:       "enableDefault.toml",
			wantRulesCount: defaultRulesCount,
			wantEnabledRules: []string{
				"noVarDeclaration",  // default rule
				"usePackageComments", // default rule
			},
			wantDisabledRules: []string{
				"noDeepExit", // non-default rule
			},
		},
		"enableDefaultRules with 2 disabled rules": {
			confPath:       "enableDefaultBut2.toml",
			wantRulesCount: defaultRulesCount - 2,
			wantEnabledRules: []string{
				"noVarDeclaration",  // default rule
				"usePackageComments", // default rule
			},
			wantDisabledRules: []string{
				"useExportedComment",          // default rule
				"lint/style/useIndentErrorFlow", // default rule
			},
		},
		"enableDefaultRules plus 1 non-default rule": {
			confPath:       "enableDefaultPlus1.toml",
			wantRulesCount: defaultRulesCount + 1,
			wantEnabledRules: []string{
				"noVarDeclaration",  // default rule
				"usePackageComments", // default rule
				"noExcessiveArguments", // non-default rule
			},
			wantDisabledRules: []string{
				"noDeepExit", // non-default rule
			},
		},
		"enableDefaultRules plus rule already in defaults": {
			confPath:       "enableDefaultPlusDefaultRule.toml",
			wantRulesCount: defaultRulesCount,
			wantEnabledRules: []string{
				"noVarDeclaration",  // default rule
				"usePackageComments", // default rule
				"useExportedComment",         // default rule
			},
			wantDisabledRules: []string{
				"noDeepExit", // non-default rule
			},
		},
		"enableAllRules plus rule already in all": {
			confPath:       "enableAllWithRule.toml",
			wantRulesCount: allRulesCount,
			wantEnabledRules: []string{
				"noVarDeclaration",  // default rule
				"usePackageComments", // default rule
				"noDeepExit",        // non-default rule
				"noExcessiveArguments", // non-default rule
			},
		},
		"enable 2 rules": {
			confPath:       "enable2.toml",
			wantRulesCount: 2,
			wantEnabledRules: []string{
				"useExportedComment",   // default rule
				"noExcessiveArguments", // non-default rule
			},
			wantDisabledRules: []string{
				"noVarDeclaration",  // default rule
				"usePackageComments", // default rule
				"noDeepExit",        // non-default rule
			},
		},
		"var-naming configure error": {
			confPath: "varNamingConfigureError.toml",
			wantErr:  `cannot configure rule: "lint/style/useVarNaming": invalid argument to the var-naming rule. Expecting a allowlist of type slice with initialisms, got string`,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			cfg, err := config.GetConfig(filepath.Join("testdata", tc.confPath))
			if err != nil {
				t.Fatalf("Unexpected error while loading conf: %v", err)
			}
			rules, err := config.GetLintingRules(cfg, []lint.Rule{})
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Fatalf("Expected error %q, got %q", tc.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error\n\t%v", err)
			}

			ruleNames := make([]string, len(rules))
			for i, rule := range rules {
				ruleNames[i] = rule.Name()
			}
			slices.Sort(ruleNames)

			if len(rules) != tc.wantRulesCount {
				t.Errorf("Expected %v enabled linting rules got: %v. Got rules: %v", tc.wantRulesCount, len(rules), ruleNames)
			}
			for _, wantEnabledRule := range tc.wantEnabledRules {
				if !slices.Contains(ruleNames, wantEnabledRule) {
					t.Errorf("Expected enabled rule %q not found. Got enabled rules: %v", wantEnabledRule, ruleNames)
				}
			}
			for _, wantDisabledRule := range tc.wantDisabledRules {
				if slices.Contains(ruleNames, wantDisabledRule) {
					t.Errorf("Expected disabled rule %q not found. Got enabled rules: %v", wantDisabledRule, ruleNames)
				}
			}
		})
	}
}

func TestGetGlobalSeverity(t *testing.T) {
	tt := map[string]struct {
		confPath               string
		wantGlobalSeverity     string
		particularRule         lint.Rule
		wantParticularSeverity string
	}{
		"enable 2 rules with one specific severity": {
			confPath:               "testdata/enable2OneSpecificSeverity.toml",
			wantGlobalSeverity:     "warning",
			particularRule:         &no_excessive_arguments.NoExcessiveArgumentsRule{},
			wantParticularSeverity: "error",
		},
		"enableAllRules with one specific severity": {
			confPath:               "testdata/enableAllOneSpecificSeverity.toml",
			wantGlobalSeverity:     "error",
			particularRule:         &no_deep_exit.DeepExitRule{},
			wantParticularSeverity: "warning",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			cfg, err := config.GetConfig(tc.confPath)
			if err != nil {
				t.Fatalf("Unexpected error while loading conf: %v", err)
			}
			rules, err := config.GetLintingRules(cfg, []lint.Rule{})
			if err != nil {
				t.Fatalf("Unexpected error while loading conf: %v", err)
			}
			for _, r := range rules {
				fullName := lint.FullRuleName(r)
				ruleCfg := cfg.Rules[fullName]
				ruleSeverity := string(ruleCfg.Severity)
				switch r.Name() {
				case tc.particularRule.Name():
					if tc.wantParticularSeverity != ruleSeverity {
						t.Fatalf("Expected Severity %v for rule %v, got %v", tc.wantParticularSeverity, fullName, ruleSeverity)
					}
				default:
					if tc.wantGlobalSeverity != ruleSeverity {
						t.Fatalf("Expected Severity %v for rule %v, got %v", tc.wantGlobalSeverity, fullName, ruleSeverity)
					}
				}
			}
		})
	}
}

func TestGetFormatter(t *testing.T) {
	t.Run("default formatter", func(t *testing.T) {
		formatter, err := config.GetFormatter("")
		if err != nil {
			t.Fatalf("Unexpected error %q", err)
		}
		if formatter == nil || formatter.Name() != "default" {
			t.Errorf("Expected formatter %q, got %v", "default", formatter)
		}
	})
	t.Run("unknown formatter", func(t *testing.T) {
		_, err := config.GetFormatter("unknown")
		if err == nil || err.Error() != "unknown formatter unknown" {
			t.Errorf("Expected error %q, got: %q", "unknown formatter unknown", err)
		}
	})
	t.Run("checkstyle formatter", func(t *testing.T) {
		formatter, err := config.GetFormatter("checkstyle")
		if err != nil {
			t.Fatalf("Unexpected error: %q", err)
		}
		if formatter == nil || formatter.Name() != "checkstyle" {
			t.Errorf("Expected formatter %q, got %v", "checkstyle", formatter)
		}
	})
}

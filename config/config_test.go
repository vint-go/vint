package config_test

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	goversion "github.com/hashicorp/go-version"

	"github.com/strowk/vint/config"
	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
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
						"blank-imports": {
							Severity: lint.SeverityWarning,
						},
						"context-as-argument": {
							Severity: lint.SeverityWarning,
						},
						"context-keys-type": {
							Severity: lint.SeverityWarning,
						},
						"dot-imports": {
							Severity: lint.SeverityWarning,
						},
						"empty-block": {
							Severity: lint.SeverityWarning,
						},
						"error-naming": {
							Severity: lint.SeverityWarning,
						},
						"error-return": {
							Severity: lint.SeverityWarning,
						},
						"error-strings": {
							Severity: lint.SeverityWarning,
						},
						"errorf": {
							Severity: lint.SeverityWarning,
						},
						"exported": {
							Severity: lint.SeverityWarning,
						},
						"increment-decrement": {
							Severity: lint.SeverityWarning,
						},
						"indent-error-flow": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noUnclosedBodies": {
							Severity: lint.SeverityWarning,
						},
						// "lint/complexity/noDuplicateCode": {
						// 	Severity: lint.SeverityWarning,
						// },
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
						"lint/correctness/noDirectErrorComparison": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noDynamicErrors": {
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
						"lint/correctness/noNilFuncComparison": {
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
						"lint/correctness/noSlogKeyValueMismatch": {
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
						"lint/correctness/noMalformedTestFunction": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noNonPointerUnmarshal": {
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
						// "lint/correctness/noMisspelledWords": {
						// 	Severity: lint.SeverityWarning,
						// },
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
						// "lint/correctness/noUnusedField": {
						// 	Severity: lint.SeverityWarning,
						// },
						"lint/suspicious/noMismatchedAppendAssign": {
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
						"lint/correctness/noDuplicateArgument": {
							Severity: lint.SeverityWarning,
						},
						"lint/suspicious/noDuplicateBranchBody": {
							Severity: lint.SeverityWarning,
						},
						"lint/correctness/noDuplicateCase": {
							Severity: lint.SeverityWarning,
						},
						"package-comments": {
							Severity: lint.SeverityWarning,
						},
						"range": {
							Severity: lint.SeverityWarning,
						},
						"receiver-naming": {
							Severity: lint.SeverityWarning,
						},
						"redefines-builtin-id": {
							Severity: lint.SeverityWarning,
						},
						"superfluous-else": {
							Severity: lint.SeverityWarning,
						},
						"time-naming": {
							Severity: lint.SeverityWarning,
						},
						"unexported-return": {
							Severity: lint.SeverityWarning,
						},
						"unreachable-code": {
							Severity: lint.SeverityWarning,
						},
						"unused-parameter": {
							Severity: lint.SeverityWarning,
						},
						"var-declaration": {
							Severity: lint.SeverityWarning,
						},
						"var-naming": {
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
						"blank-imports":                     {},
						"context-as-argument":               {},
						"context-keys-type":                 {},
						"dot-imports":                       {},
						"empty-block":                       {},
						"error-naming":                      {},
						"error-return":                      {},
						"error-strings":                     {},
						"errorf":                            {},
						"exported":                          {},
						"increment-decrement":               {},
						"indent-error-flow":                 {},
						"lint/correctness/noUnclosedBodies": {},
						// "lint/complexity/noDuplicateCode":       {},
						"lint/complexity/noExcessiveStatements":          {},
						"lint/complexity/noHighCyclomaticComplexity":     {},
						"lint/complexity/noLongFunctions":                {},
						"lint/correctness/noDeniedImport":                {},
						"lint/correctness/noDirectErrorComparison":       {},
						"lint/correctness/noDynamicErrors":               {},
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
						"lint/style/noInitFunction":                      {},
						"lint/style/noLeadingBlankLine":                  {},
						"lint/style/noTrailingBlankLine":                 {},
						"lint/style/noRepeatedStrings":                   {},
						"lint/style/noUnnecessaryConversion":             {},
						"lint/style/noUnnecessaryLoopVarCopy":            {},
						"lint/style/usePrintfSuffix":                     {},
						"lint/correctness/noSingleArgAppend":             {},
						"lint/correctness/noAsmDeclMismatch":             {},
						"lint/correctness/noSelfAssignment":              {},
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
						"lint/correctness/noInvalidErrorsAs":             {},
						"lint/correctness/noFramePointerClobber":         {},
						"lint/correctness/useJoinHostPort":               {},
						"lint/correctness/noImpossibleInterfaceAssert":   {},
						"lint/correctness/noLoopClosureCapture":          {},
						"lint/correctness/noLostCancel":                  {},
						"lint/correctness/noWaitGroupMisuse":             {},
						"lint/correctness/noNilFuncComparison":           {},
						"lint/correctness/noPrintfFormatMismatch":        {},
						"lint/correctness/noExcessiveShift":              {},
						"lint/correctness/noUnbufferedSignalChannel":     {},
						"lint/correctness/noSlogKeyValueMismatch":        {},
						"lint/correctness/noStdMethodSignatureMismatch":  {},
						"lint/correctness/noStdlibVersionMismatch":       {},
						"lint/correctness/noStringIntConversion":         {},
						"lint/correctness/noTestFatalInGoroutine":        {},
						"lint/correctness/noMalformedTestFunction":       {},
						"lint/correctness/noNonPointerUnmarshal":         {},
						"lint/correctness/noUnreachableCode":             {},
						"lint/correctness/noInvalidUnsafePointer":        {},
						"lint/correctness/noUnusedFunctionResult":        {},
						"lint/correctness/noUnusedWrite":                 {},
						"lint/correctness/noIneffectualAssignment":       {},
						"lint/style/noLineTooLong":                       {},
						"lint/style/noMagicNumberInArgument":             {},
						"lint/style/noMagicNumberInAssignment":           {},
						"lint/style/noMagicNumberInCase":                 {},
						"lint/style/noMagicNumberInCondition":            {},
						"lint/style/noMagicNumberInOperation":            {},
						"lint/style/noMagicNumberInReturn":               {},
						// "lint/correctness/noMisspelledWords":             {},
						"lint/style/noNakedReturn":                 {},
						"lint/correctness/noTlsConnHandshake":      {},
						"lint/correctness/noTlsDial":               {},
						"lint/correctness/noTlsDialWithDialer":     {},
						"lint/correctness/noNetDial":               {},
						"lint/correctness/noNetDialTimeout":        {},
						"lint/correctness/noNetListen":             {},
						"lint/correctness/noNetListenPacket":       {},
						"lint/correctness/noNetLookupCname":        {},
						"lint/correctness/noNetLookupHost":         {},
						"lint/correctness/noNetLookupIp":           {},
						"lint/correctness/noNetLookupMx":           {},
						"lint/correctness/noNetLookupPort":         {},
						"lint/correctness/noNetLookupSrv":          {},
						"lint/correctness/noNetLookupNs":           {},
						"lint/correctness/noNetLookupTxt":          {},
						"lint/correctness/noNetLookupAddr":         {},
						"lint/correctness/noSqlDbBegin":            {},
						"lint/correctness/noSqlDbExec":             {},
						"lint/correctness/noSqlDbPing":             {},
						"lint/correctness/noSqlDbPrepare":          {},
						"lint/correctness/noSqlDbQuery":            {},
						"lint/correctness/noSqlDbQueryRow":         {},
						"lint/correctness/noSqlTxExec":             {},
						"lint/correctness/noSqlTxPrepare":          {},
						"lint/correctness/noSqlTxQuery":            {},
						"lint/correctness/noSqlTxQueryRow":         {},
						"lint/correctness/noSqlTxStmt":             {},
						"lint/correctness/noSqlStmtExec":           {},
						"lint/correctness/noSqlStmtQuery":          {},
						"lint/correctness/noSqlStmtQueryRow":       {},
						"lint/correctness/noExecCommand":           {},
						"lint/style/noNolintNonMachineReadable":    {},
						"lint/correctness/noNolintParseError":      {},
						"lint/style/noNolintWithoutSpecificLinter": {},
						"lint/style/noNolintWithoutExplanation":    {},
						"lint/suspicious/noUnusedNolint":           {},
						"lint/suspicious/noUnusedParameter":        {},
						"lint/suspicious/noConstantParameter":      {},
						"lint/suspicious/noConstantResult":         {},
						"lint/correctness/noUnusedFunction":        {},
						"lint/correctness/noUnusedType":            {},
						"lint/correctness/noUnusedVariable":        {},
						"lint/correctness/noUnusedConstant":        {},
						// "lint/correctness/noUnusedField":           {},
						"lint/suspicious/noMismatchedAppendAssign": {},
						"lint/performance/useCombinedAppend":       {},
						"lint/suspicious/noSwappedArguments":       {},
						"lint/style/useAssignmentOperator":         {},
						"lint/correctness/noNoopFunctionCall":      {},
						"lint/correctness/noImpossibleCondition":   {},
						"lint/style/noCapitalizedLocal":            {},
						"lint/correctness/noUnreachableTypeCase":   {},
						"lint/style/noMisplacedDefaultCase":        {},
						"lint/style/useStandardCodegenComment":     {},
						"lint/style/useStandardDeprecationComment": {},
						"lint/style/useCommentSpacing":             {},
						"lint/correctness/noDuplicateArgument":     {},
						"lint/suspicious/noDuplicateBranchBody":    {},
						"lint/correctness/noDuplicateCase":         {},
						"package-comments":                         {},
						"range":                                    {},
						"receiver-naming":                          {},
						"redefines-builtin-id":                     {},
						"superfluous-else":                         {},
						"time-naming":                              {},
						"unexported-return":                        {},
						"unreachable-code":                         {},
						"unused-parameter":                         {},
						"var-declaration":                          {},
						"var-naming":                               {},
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
						"argument-limit": {
							Severity: lint.SeverityWarning,
							Exclude:  []string{"excluded/file.go"},
							Arguments: lint.Arguments{
								[]any{4},
							},
						},
						"blank-imports": {
							Disabled: true,
							Severity: lint.SeverityError,
						},
						"context-as-argument": {
							Severity: lint.SeverityError,
						},
						"context-keys-type": {
							Severity: lint.SeverityError,
						},
						"dot-imports": {
							Severity: lint.SeverityError,
						},
						"empty-block": {
							Severity: lint.SeverityError,
						},
						"error-naming": {
							Severity: lint.SeverityError,
						},
						"error-return": {
							Severity: lint.SeverityError,
						},
						"error-strings": {
							Severity: lint.SeverityError,
						},
						"errorf": {
							Severity: lint.SeverityError,
						},
						"exported": {
							Severity: lint.SeverityError,
							Arguments: lint.Arguments{
								"check-private-receivers", "disable-stuttering-check",
							},
							Exclude: []string{"excluded/file-exported.go"},
						},
						"increment-decrement": {
							Severity: lint.SeverityError,
						},
						"indent-error-flow": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noUnclosedBodies": {
							Severity: lint.SeverityError,
						},
						// "lint/complexity/noDuplicateCode": {
						// 	Severity: lint.SeverityError,
						// },
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
						"lint/correctness/noDirectErrorComparison": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noDynamicErrors": {
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
						"lint/correctness/noNilFuncComparison": {
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
						"lint/correctness/noSlogKeyValueMismatch": {
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
						"lint/correctness/noMalformedTestFunction": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noNonPointerUnmarshal": {
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
						// "lint/correctness/noMisspelledWords": {
						// 	Severity: lint.SeverityError,
						// },
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
						// "lint/correctness/noUnusedField": {
						// 	Severity: lint.SeverityError,
						// },
						"lint/suspicious/noMismatchedAppendAssign": {
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
						"lint/correctness/noDuplicateArgument": {
							Severity: lint.SeverityError,
						},
						"lint/suspicious/noDuplicateBranchBody": {
							Severity: lint.SeverityError,
						},
						"lint/correctness/noDuplicateCase": {
							Severity: lint.SeverityError,
						},
						"package-comments": {
							Severity: lint.SeverityError,
						},
						"range": {
							Severity: lint.SeverityError,
						},
						"receiver-naming": {
							Severity: lint.SeverityError,
						},
						"redefines-builtin-id": {
							Severity: lint.SeverityError,
						},
						"superfluous-else": {
							Severity: lint.SeverityError,
						},
						"time-naming": {
							Severity: lint.SeverityError,
						},
						"unexported-return": {
							Severity: lint.SeverityError,
						},
						"unreachable-code": {
							Severity: lint.SeverityError,
						},
						"unused-parameter": {
							Severity: lint.SeverityError,
						},
						"var-declaration": {
							Severity: lint.SeverityError,
						},
						"var-naming": {
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
				wantError: "error in config of rule [var-naming]",
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
		defaultRulesCount = 208
		// len of allRules: update this when adding new rules
		allRulesCount = 315
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
				"var-declaration",  // default rule
				"package-comments", // default rule
				"deep-exit",        // non-default rule
			},
		},
		"enableAllRules without disabled rules": {
			confPath:       "enableAll.toml",
			wantRulesCount: allRulesCount,
			wantEnabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
				"deep-exit",        // non-default rule
			},
		},
		"enableAllRules with 2 disabled rules": {
			confPath:       "enableAllBut2.toml",
			wantRulesCount: allRulesCount - 2,
			wantEnabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
				"deep-exit",        // non-default rule
			},
			wantDisabledRules: []string{
				"exported",   // default rule
				"cyclomatic", // non-default rule
			},
		},
		"enableDefaultRules without disabled rules": {
			confPath:       "enableDefault.toml",
			wantRulesCount: defaultRulesCount,
			wantEnabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
			},
			wantDisabledRules: []string{
				"deep-exit", // non-default rule
			},
		},
		"enableDefaultRules with 2 disabled rules": {
			confPath:       "enableDefaultBut2.toml",
			wantRulesCount: defaultRulesCount - 2,
			wantEnabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
			},
			wantDisabledRules: []string{
				"exported",          // default rule
				"indent-error-flow", // default rule
			},
		},
		"enableDefaultRules plus 1 non-default rule": {
			confPath:       "enableDefaultPlus1.toml",
			wantRulesCount: defaultRulesCount + 1,
			wantEnabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
				"cyclomatic",       // non-default rule
			},
			wantDisabledRules: []string{
				"deep-exit", // non-default rule
			},
		},
		"enableDefaultRules plus rule already in defaults": {
			confPath:       "enableDefaultPlusDefaultRule.toml",
			wantRulesCount: defaultRulesCount,
			wantEnabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
				"exported",         // default rule
			},
			wantDisabledRules: []string{
				"deep-exit", // non-default rule
			},
		},
		"enableAllRules plus rule already in all": {
			confPath:       "enableAllWithRule.toml",
			wantRulesCount: allRulesCount,
			wantEnabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
				"deep-exit",        // non-default rule
				"cyclomatic",       // non-default rule
			},
		},
		"enable 2 rules": {
			confPath:       "enable2.toml",
			wantRulesCount: 2,
			wantEnabledRules: []string{
				"exported",   // default rule
				"cyclomatic", // non-default rule
			},
			wantDisabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
				"deep-exit",        // non-default rule
			},
		},
		"enable imports-blocklist rule": {
			confPath:       "issue-969.toml",
			wantRulesCount: 1,
			wantEnabledRules: []string{
				"imports-blocklist", // non-default renamed rule
			},
			wantDisabledRules: []string{
				"imports-blacklist", // non-default deprecated rule name
			},
		},
		"var-naming configure error": {
			confPath: "varNamingConfigureError.toml",
			wantErr:  `cannot configure rule: "var-naming": invalid argument to the var-naming rule. Expecting a allowlist of type slice with initialisms, got string`,
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
			particularRule:         &rule.CyclomaticRule{},
			wantParticularSeverity: "error",
		},
		"enableAllRules with one specific severity": {
			confPath:               "testdata/enableAllOneSpecificSeverity.toml",
			wantGlobalSeverity:     "error",
			particularRule:         &rule.DeepExitRule{},
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

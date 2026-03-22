package gosec

import "github.com/strowk/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the gosec golangci-lint linter.
// gosec inspects source code for security problems by scanning the Go AST
// and SSA code representation. It supports includes/excludes to control
// which rules are active, plus per-rule configuration.
type Migrator struct{}

func (*Migrator) Name() string {
	return "gosec"
}

// allRules maps gosec rule IDs to their vint rule paths.
// Only rules that exist in the vint rules directory are included.
var allRules = map[string]string{
	"G101": "lint/security/noHardcodedCredentials",
	"G102": "lint/security/noBindToAllInterfaces",
	"G103": "lint/security/noUnsafePackage",
	"G104": "lint/correctness/noUncheckedError", // subsumed: gosec unchecked errors
	"G106": "lint/security/noInsecureHostKeyCallback",
	"G107": "lint/security/noSsrfViaVariable",
	"G108": "lint/security/noExposedPprof",
	"G109": "lint/security/noAtoiOverflow",
	"G110": "lint/security/noUnboundedDecompression",
	"G111": "lint/security/noFilesystemRootServing",
	"G112": "lint/security/noMissingReadHeaderTimeout",
	"G113": "lint/security/noHttpRequestSmuggling",
	"G114": "lint/security/noServeWithoutTimeout",
	"G115": "lint/security/noIntegerOverflowConversion",
	"G116": "lint/security/noTrojanSourceBidi",
	"G117": "lint/security/noSecretInSerialization",
	"G118": "lint/security/noContextPropagationFailure",
	"G119": "lint/security/noUnsafeRedirectPolicy",
	"G120": "lint/security/noUnboundedFormParsing",
	"G121": "lint/security/noUnsafeCorsBypass",
	"G122": "lint/security/noFilesystemToctou",
	"G123": "lint/security/noTlsSessionResumptionBypass",
	"G124": "lint/security/noInsecureCookie",
	"G201": "lint/security/noSqlFormatString",
	"G202": "lint/security/noSqlConcatenation",
	"G203": "lint/security/noUnescapedHtmlTemplate",
	"G204": "lint/security/noVariableCommandExecution",
	"G301": "lint/security/noPermissiveDirectoryPermissions",
	"G302": "lint/security/noPermissiveFilePermissions",
	"G303": "lint/security/noPredictableTempFile",
	"G304": "lint/security/noFileInclusionViaVariable",
	"G305": "lint/security/noZipSlip",
	"G306": "lint/security/noPermissiveWriteFilePermissions",
	"G307": "lint/security/noPermissiveOsCreate",
	"G401": "lint/security/noWeakCryptoHash",
	"G402": "lint/security/noInsecureTlsConfig",
	"G403": "lint/security/noShortRsaKey",
	"G404": "lint/security/noInsecureRandom",
	"G405": "lint/security/noWeakEncryptionAlgorithm",
	"G406": "lint/security/noDeprecatedHashFunction",
	"G407": "lint/security/noHardcodedIv",
	"G408": "lint/security/noSshAuthBypass",
	"G501": "lint/security/noWeakCryptoHash",          // subsumed: import crypto/md5
	"G502": "lint/security/noWeakEncryptionAlgorithm", // subsumed: import crypto/des
	"G503": "lint/security/noWeakEncryptionAlgorithm", // subsumed: import crypto/rc4
	"G504": "lint/security/noCgiImport",
	"G505": "lint/security/noWeakCryptoHash",          // subsumed: import crypto/sha1
	"G506": "lint/security/noDeprecatedHashFunction",  // subsumed: import x/crypto/md4
	"G507": "lint/security/noDeprecatedHashFunction",  // subsumed: import x/crypto/ripemd160
	"G601": "lint/security/noRangeVariableAlias",
	"G602": "lint/security/noSliceBoundsOutOfRange",
	"G701": "lint/security/noSqlInjectionTaint",
	"G702": "lint/security/noCommandInjectionTaint",
	"G703": "lint/security/noPathTraversalTaint",
	"G704": "lint/security/noSsrfTaint",
	"G705": "lint/security/noXssTaint",
	"G706": "lint/security/noLogInjectionTaint",
	"G707": "lint/security/noSmtpInjectionTaint",
	"G708": "lint/security/noTemplateInjection",
	"G709": "lint/security/noUnsafeDeserialization",
}

func (*Migrator) MigrateConfig(settings map[string]any) (map[string]migrate.VintRuleConfig, error) {
	configs := make(map[string]migrate.VintRuleConfig)

	// Determine which gosec rule IDs are active based on includes/excludes.
	activeIDs := resolveActiveRules(settings)

	// Enable each active rule in vint config.
	for _, id := range sortedKeys(activeIDs) {
		vintPath, ok := allRules[id]
		if !ok {
			// Rule ID not mapped to a vint rule (e.g. missing rules).
			continue
		}
		configs[vintPath] = migrate.VintRuleConfig{}
	}

	return configs, nil
}

// resolveActiveRules determines which gosec rule IDs should be active
// based on the includes and excludes settings.
func resolveActiveRules(settings map[string]any) map[string]bool {
	active := make(map[string]bool)

	// Check for includes first.
	var includes []string
	if settings != nil {
		if v, ok := settings["includes"]; ok {
			includes = toStringSlice(v)
		}
	}

	if len(includes) > 0 {
		// Only include specified rules.
		for _, id := range includes {
			if _, ok := allRules[id]; ok {
				active[id] = true
			}
		}
	} else {
		// Default: all rules are active.
		for id := range allRules {
			active[id] = true
		}
	}

	// Apply excludes.
	if settings != nil {
		if v, ok := settings["excludes"]; ok {
			for _, id := range toStringSlice(v) {
				delete(active, id)
			}
		}
	}

	return active
}

// toStringSlice converts an any value (expected to be []any with string elements)
// to a []string.
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

// sortedKeys returns the keys of a map in sorted order.
func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Sort for deterministic output.
	sortStrings(keys)
	return keys
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

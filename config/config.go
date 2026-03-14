// Package config implements revive's configuration data structures and related methods.
package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"

	"github.com/strowk/vint/formatter"
	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
	"github.com/strowk/vint/rules/no_blank_error_assignment"
	"github.com/strowk/vint/rules/no_denied_import"
	"github.com/strowk/vint/rules/no_direct_error_comparison"
	"github.com/strowk/vint/rules/no_duplicate_code"
	"github.com/strowk/vint/rules/no_dynamic_errors"
	"github.com/strowk/vint/rules/no_excessive_statements"
	"github.com/strowk/vint/rules/no_file_scoped_denied_import"
	"github.com/strowk/vint/rules/no_high_cyclomatic_complexity"
	"github.com/strowk/vint/rules/no_long_functions"
	"github.com/strowk/vint/rules/no_range_variable_alias"
	"github.com/strowk/vint/rules/no_slice_bounds_out_of_range"
	"github.com/strowk/vint/rules/no_space_in_directive"
	"github.com/strowk/vint/rules/no_unallowed_import"
	"github.com/strowk/vint/rules/no_unrecognized_directive"
	"github.com/strowk/vint/rules/no_atoi_overflow"
	"github.com/strowk/vint/rules/no_cgi_import"
	"github.com/strowk/vint/rules/no_bind_to_all_interfaces"
	"github.com/strowk/vint/rules/no_exposed_pprof"
	"github.com/strowk/vint/rules/no_filesystem_root_serving"
	"github.com/strowk/vint/rules/no_hardcoded_credentials"
	"github.com/strowk/vint/rules/no_hardcoded_iv"
	"github.com/strowk/vint/rules/no_http_request_smuggling"
	"github.com/strowk/vint/rules/no_insecure_cookie"
	"github.com/strowk/vint/rules/no_insecure_host_key_callback"
	"github.com/strowk/vint/rules/no_integer_overflow_conversion"
	"github.com/strowk/vint/rules/no_missing_read_header_timeout"
	"github.com/strowk/vint/rules/no_serve_without_timeout"
	"github.com/strowk/vint/rules/no_ssrf_via_variable"
	"github.com/strowk/vint/rules/no_trojan_source_bidi"
	"github.com/strowk/vint/rules/no_unchecked_error"
	"github.com/strowk/vint/rules/no_unbounded_decompression"
	"github.com/strowk/vint/rules/no_unchecked_type_assertion"
	"github.com/strowk/vint/rules/no_unclosed_bodies"
	"github.com/strowk/vint/rules/no_unsafe_package"
	"github.com/strowk/vint/rules/no_command_injection_taint"
	"github.com/strowk/vint/rules/no_context_propagation_failure"
	"github.com/strowk/vint/rules/no_deprecated_hash_function"
	"github.com/strowk/vint/rules/no_file_inclusion_via_variable"
	"github.com/strowk/vint/rules/no_filesystem_toctou"
	"github.com/strowk/vint/rules/no_insecure_random"
	"github.com/strowk/vint/rules/no_insecure_tls_config"
	"github.com/strowk/vint/rules/no_log_injection_taint"
	"github.com/strowk/vint/rules/no_path_traversal_taint"
	"github.com/strowk/vint/rules/no_permissive_directory_permissions"
	"github.com/strowk/vint/rules/no_permissive_file_permissions"
	"github.com/strowk/vint/rules/no_permissive_os_create"
	"github.com/strowk/vint/rules/no_permissive_write_file_permissions"
	"github.com/strowk/vint/rules/no_predictable_temp_file"
	"github.com/strowk/vint/rules/no_secret_in_serialization"
	"github.com/strowk/vint/rules/no_short_rsa_key"
	"github.com/strowk/vint/rules/no_smtp_injection_taint"
	"github.com/strowk/vint/rules/no_sql_concatenation"
	"github.com/strowk/vint/rules/no_sql_format_string"
	"github.com/strowk/vint/rules/no_sql_injection_taint"
	"github.com/strowk/vint/rules/no_ssh_auth_bypass"
)

var defaultRules = []lint.Rule{
	&rule.VarDeclarationsRule{},
	&rule.PackageCommentsRule{},
	&rule.DotImportsRule{},
	&rule.BlankImportsRule{},
	&rule.ExportedRule{},
	&rule.VarNamingRule{},
	&rule.IndentErrorFlowRule{},
	&rule.RangeRule{},
	&rule.ErrorfRule{},
	&rule.ErrorNamingRule{},
	&rule.ErrorStringsRule{},
	&rule.ReceiverNamingRule{},
	&rule.IncrementDecrementRule{},
	&rule.ErrorReturnRule{},
	&rule.UnexportedReturnRule{},
	&rule.TimeNamingRule{},
	&rule.ContextKeysType{},
	&rule.ContextAsArgumentRule{},
	&rule.EmptyBlockRule{},
	&rule.SuperfluousElseRule{},
	&rule.UnusedParamRule{},
	&rule.UnreachableCodeRule{},
	&rule.RedefinesBuiltinIDRule{},
	&no_unclosed_bodies.NoUnclosedBodiesRule{},
	&no_duplicate_code.NoDuplicateCodeRule{},
	&no_excessive_statements.NoExcessiveStatementsRule{},
	&no_high_cyclomatic_complexity.NoHighCyclomaticComplexityRule{},
	&no_long_functions.NoLongFunctionsRule{},
	&no_denied_import.NoDeniedImportRule{},
	&no_direct_error_comparison.NoDirectErrorComparisonRule{},
	&no_dynamic_errors.NoDynamicErrorsRule{},
	&no_file_scoped_denied_import.NoFileScopedDeniedImportRule{},
	&no_space_in_directive.NoSpaceInDirectiveRule{},
	&no_unallowed_import.NoUnallowedImportRule{},
	&no_unrecognized_directive.NoUnrecognizedDirectiveRule{},
	&no_unchecked_error.NoUncheckedErrorRule{},
	&no_range_variable_alias.NoRangeVariableAliasRule{},
	&no_slice_bounds_out_of_range.NoSliceBoundsOutOfRangeRule{},
	&no_atoi_overflow.NoAtoiOverflowRule{},
	&no_bind_to_all_interfaces.NoBindToAllInterfacesRule{},
	&no_exposed_pprof.NoExposedPprofRule{},
	&no_filesystem_root_serving.NoFilesystemRootServingRule{},
	&no_hardcoded_credentials.NoHardcodedCredentialsRule{},
	&no_hardcoded_iv.NoHardcodedIvRule{},
	&no_http_request_smuggling.NoHttpRequestSmugglingRule{},
	&no_insecure_cookie.NoInsecureCookieRule{},
	&no_insecure_host_key_callback.NoInsecureHostKeyCallbackRule{},
	&no_integer_overflow_conversion.NoIntegerOverflowConversionRule{},
	&no_missing_read_header_timeout.NoMissingReadHeaderTimeoutRule{},
	&no_serve_without_timeout.NoServeWithoutTimeoutRule{},
	&no_ssrf_via_variable.NoSsrfViaVariableRule{},
	&no_trojan_source_bidi.NoTrojanSourceBidiRule{},
	&no_unbounded_decompression.NoUnboundedDecompressionRule{},
	&no_unsafe_package.NoUnsafePackageRule{},
	&no_cgi_import.NoCgiImportRule{},
	&no_command_injection_taint.NoCommandInjectionTaintRule{},
	&no_context_propagation_failure.NoContextPropagationFailureRule{},
	&no_deprecated_hash_function.NoDeprecatedHashFunctionRule{},
	&no_file_inclusion_via_variable.NoFileInclusionViaVariableRule{},
	&no_filesystem_toctou.NoFilesystemToctouRule{},
	&no_insecure_random.NoInsecureRandomRule{},
	&no_insecure_tls_config.NoInsecureTlsConfigRule{},
	&no_log_injection_taint.NoLogInjectionTaintRule{},
	&no_path_traversal_taint.NoPathTraversalTaintRule{},
	&no_permissive_directory_permissions.NoPermissiveDirectoryPermissionsRule{},
	&no_permissive_file_permissions.NoPermissiveFilePermissionsRule{},
	&no_permissive_os_create.NoPermissiveOsCreateRule{},
	&no_permissive_write_file_permissions.NoPermissiveWriteFilePermissionsRule{},
	&no_predictable_temp_file.NoPredictableTempFileRule{},
	&no_secret_in_serialization.NoSecretInSerializationRule{},
	&no_short_rsa_key.NoShortRsaKeyRule{},
	&no_smtp_injection_taint.NoSmtpInjectionTaintRule{},
	&no_sql_concatenation.NoSqlConcatenationRule{},
	&no_sql_format_string.NoSqlFormatStringRule{},
	&no_sql_injection_taint.NoSqlInjectionTaintRule{},
	&no_ssh_auth_bypass.NoSshAuthBypassRule{},
}

var allRules = append([]lint.Rule{
	&rule.ArgumentsLimitRule{},
	&rule.CyclomaticRule{},
	&rule.FileHeaderRule{},
	&rule.ConfusingNamingRule{},
	&rule.GetReturnRule{},
	&rule.ModifiesParamRule{},
	&rule.ConfusingResultsRule{},
	&rule.DeepExitRule{},
	&rule.AddConstantRule{},
	&rule.FlagParamRule{},
	&rule.UnnecessaryStmtRule{},
	&rule.StructTagRule{},
	&rule.ModifiesValRecRule{},
	&rule.ConstantLogicalExprRule{},
	&rule.BoolLiteralRule{},
	&rule.ImportsBlocklistRule{},
	&rule.FunctionResultsLimitRule{},
	&rule.MaxPublicStructsRule{},
	&rule.RangeValInClosureRule{},
	&rule.RangeValAddress{},
	&rule.WaitGroupByValueRule{},
	&rule.AtomicRule{},
	&rule.EmptyLinesRule{},
	&rule.LineLengthLimitRule{},
	&rule.CallToGCRule{},
	&rule.DuplicatedImportsRule{},
	&rule.ImportShadowingRule{},
	&rule.BareReturnRule{},
	&rule.UnusedReceiverRule{},
	&rule.UnhandledErrorRule{},
	&rule.CognitiveComplexityRule{},
	&rule.StringOfIntRule{},
	&rule.StringFormatRule{},
	&rule.EarlyReturnRule{},
	&rule.UnconditionalRecursionRule{},
	&rule.IdenticalBranchesRule{},
	&rule.DeferRule{},
	&rule.UnexportedNamingRule{},
	&rule.FunctionLength{},
	&rule.NestedStructs{},
	&rule.UselessBreak{},
	&rule.UncheckedTypeAssertionRule{},
	&rule.TimeEqualRule{},
	&rule.TimeDateRule{},
	&rule.BannedCharsRule{},
	&rule.OptimizeOperandsOrderRule{},
	&rule.UseAnyRule{},
	&rule.DataRaceRule{},
	&rule.CommentSpacingsRule{},
	&rule.IfReturnRule{},
	&rule.RedundantImportAlias{},
	&rule.ImportAliasNamingRule{},
	&rule.EnforceMapStyleRule{},
	&rule.EnforceRepeatedArgTypeStyleRule{},
	&rule.EnforceSliceStyleRule{},
	&rule.MaxControlNestingRule{},
	&rule.CommentsDensityRule{},
	&rule.FileLengthLimitRule{},
	&rule.FilenameFormatRule{},
	&rule.RedundantBuildTagRule{},
	&rule.UseErrorsNewRule{},
	&rule.RedundantTestMainExitRule{},
	&rule.UnnecessaryFormatRule{},
	&rule.UseFmtPrintRule{},
	&rule.EnforceSwitchStyleRule{},
	&rule.IdenticalSwitchConditionsRule{},
	&rule.IdenticalIfElseIfConditionsRule{},
	&rule.IdenticalIfElseIfBranchesRule{},
	&rule.IdenticalSwitchBranchesRule{},
	&rule.UselessFallthroughRule{},
	&rule.PackageDirectoryMismatchRule{},
	&rule.UseWaitGroupGoRule{},
	&rule.UnsecureURLSchemeRule{},
	&rule.InefficientMapLookupRule{},
	&rule.ForbiddenCallInWgGoRule{},
	&rule.UnnecessaryIfRule{},
	&rule.EpochNamingRule{},
	&rule.UseSlicesSort{},
	&rule.PackageNamingRule{},
	&no_blank_error_assignment.NoBlankErrorAssignmentRule{},
	&no_unchecked_type_assertion.NoUncheckedTypeAssertionRule{},
}, defaultRules...)

// allFormatters is a list of all available formatters to output the linting results.
// Keep the list sorted and in sync with available formatters in README.md.
var allFormatters = []lint.Formatter{
	&formatter.Checkstyle{},
	&formatter.Default{},
	&formatter.Friendly{},
	&formatter.JSON{},
	&formatter.NDJSON{},
	&formatter.Plain{},
	&formatter.Sarif{},
	&formatter.Stylish{},
	&formatter.Unix{},
}

func getFormatters() map[string]lint.Formatter {
	result := map[string]lint.Formatter{}
	for _, f := range allFormatters {
		result[f.Name()] = f
	}
	return result
}

// GetLintingRules yields the linting rules that must be applied by the linter.
func GetLintingRules(config *lint.Config, extraRules []lint.Rule) ([]lint.Rule, error) {
	rulesMap := map[string]lint.Rule{}
	for _, r := range allRules {
		rulesMap[lint.FullRuleName(r)] = r
	}
	for _, r := range extraRules {
		fullName := lint.FullRuleName(r)
		if _, ok := rulesMap[fullName]; ok {
			continue
		}
		rulesMap[fullName] = r
	}

	var lintingRules []lint.Rule
	for name, ruleConfig := range config.Rules {
		actualName := actualRuleName(name)
		r, ok := rulesMap[actualName]
		if !ok {
			return nil, fmt.Errorf("cannot find rule: %s", name)
		}

		if ruleConfig.Disabled {
			continue // skip disabled rules
		}

		if r, ok := r.(lint.ConfigurableRule); ok {
			if err := r.Configure(ruleConfig.Arguments); err != nil {
				return nil, fmt.Errorf("cannot configure rule: %q: %w", name, err)
			}
		}

		lintingRules = append(lintingRules, r)
	}

	return lintingRules, nil
}

func actualRuleName(name string) string {
	switch name {
	case "imports-blacklist":
		return "imports-blocklist"
	default:
		return name
	}
}

func parseConfig(data []byte, config *lint.Config) error {
	err := toml.Unmarshal(data, config)
	if err != nil {
		return fmt.Errorf("cannot parse the config file: %w", err)
	}
	for k, r := range config.Rules {
		err := r.Initialize()
		if err != nil {
			return fmt.Errorf("error in config of rule [%s] : [%w]", k, err)
		}
		config.Rules[k] = r
	}

	return nil
}

// vintYAMLConfig represents the vint.yaml configuration file structure.
type vintYAMLConfig struct {
	Settings map[string]vintYAMLRuleConfig `yaml:"settings"`
}

// vintYAMLRuleConfig represents a single rule's configuration in vint.yaml.
type vintYAMLRuleConfig struct {
	Severity  string `yaml:"severity"`
	Disabled  bool   `yaml:"disabled"`
	Arguments []any  `yaml:"arguments"`
	Exclude   []string `yaml:"exclude"`
}

// mergeVintYAML reads vint.yaml from the given path and merges its settings
// into the existing config. YAML settings override matching TOML settings.
func mergeVintYAML(path string, config *lint.Config) error {
	data, err := os.ReadFile(path) //nolint:gosec // ignore G304: potential file inclusion via variable
	if err != nil {
		return fmt.Errorf("cannot read vint.yaml: %w", err)
	}

	var yamlCfg vintYAMLConfig
	if err := yaml.Unmarshal(data, &yamlCfg); err != nil {
		return fmt.Errorf("cannot parse vint.yaml: %w", err)
	}

	if config.Rules == nil {
		config.Rules = map[string]lint.RuleConfig{}
	}

	for name, yamlRule := range yamlCfg.Settings {
		rc := lint.RuleConfig{
			Severity:  lint.Severity(yamlRule.Severity),
			Disabled:  yamlRule.Disabled,
			Arguments: yamlRule.Arguments,
			Exclude:   yamlRule.Exclude,
		}
		if err := rc.Initialize(); err != nil {
			return fmt.Errorf("error in vint.yaml config of rule [%s]: %w", name, err)
		}
		config.Rules[name] = rc
	}

	return nil
}

func validateConfig(config *lint.Config) error {
	if config.EnableAllRules && config.EnableDefaultRules {
		return errors.New("config options enableAllRules and enableDefaultRules cannot be combined")
	}
	return nil
}

func normalizeConfig(config *lint.Config) {
	if len(config.Rules) == 0 {
		config.Rules = map[string]lint.RuleConfig{}
	}

	addRules := func(config *lint.Config, rules []lint.Rule) {
		for _, r := range rules {
			ruleName := lint.FullRuleName(r)
			if _, ok := config.Rules[ruleName]; !ok {
				config.Rules[ruleName] = lint.RuleConfig{}
			}
		}
	}

	if config.EnableAllRules {
		addRules(config, allRules)
	} else if config.EnableDefaultRules {
		addRules(config, defaultRules)
	}

	severity := config.Severity
	if severity != "" {
		for k, v := range config.Rules {
			if v.Severity == "" {
				v.Severity = severity
			}
			config.Rules[k] = v
		}
		for k, v := range config.Directives {
			if v.Severity == "" {
				v.Severity = severity
			}
			config.Directives[k] = v
		}
	}
}

const defaultConfidence = 0.8

// vintYAMLPath is the hardcoded config file name for YAML configuration.
const vintYAMLPath = "vint.yaml"

// GetConfig yields the configuration.
func GetConfig(configPath string) (*lint.Config, error) {
	config := &lint.Config{}
	switch {
	case configPath != "":
		config.Confidence = defaultConfidence
		data, err := os.ReadFile(configPath) //nolint:gosec // ignore G304: potential file inclusion via variable
		if err != nil {
			return nil, errors.New("cannot read the config file")
		}
		err = parseConfig(data, config)
		if err != nil {
			return nil, err
		}

	default: // no configuration provided
		config = defaultConfig()
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	normalizeConfig(config)

	// Merge vint.yaml if it exists in CWD (overrides TOML settings)
	if _, err := os.Stat(vintYAMLPath); err == nil {
		if err := mergeVintYAML(vintYAMLPath, config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// GetFormatter yields the formatter for lint failures.
func GetFormatter(formatterName string) (lint.Formatter, error) {
	formatters := getFormatters()
	if formatterName == "" {
		return formatters["default"], nil
	}
	f, ok := formatters[formatterName]
	if !ok {
		return nil, fmt.Errorf("unknown formatter %v", formatterName)
	}
	return f, nil
}

func defaultConfig() *lint.Config {
	defaultConfig := lint.Config{
		Confidence: defaultConfidence,
		Severity:   lint.SeverityWarning,
		Rules:      map[string]lint.RuleConfig{},
	}
	for _, r := range defaultRules {
		defaultConfig.Rules[lint.FullRuleName(r)] = lint.RuleConfig{}
	}
	return &defaultConfig
}

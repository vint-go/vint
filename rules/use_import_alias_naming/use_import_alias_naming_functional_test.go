package use_import_alias_naming_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_import_alias_naming"
)

func TestUseImportAliasNaming(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_import_alias_naming", &use_import_alias_naming.ImportAliasNamingRule{})
	functional_test_helpers.TestRule(t, "use_import_alias_naming", &use_import_alias_naming.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{},
		},
	})
}

func TestUseImportAliasNaming_CustomConfig(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_import_alias_naming_custom_config", &use_import_alias_naming.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{`^[a-z]+$`},
	})
}

func TestUseImportAliasNaming_CustomConfigWithMultipleRules(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_import_alias_naming_custom_config_with_multiple_values", &use_import_alias_naming.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"allowRegex": `^[a-z][a-z0-9]*$`,
				"denyRegex":  `^((v\d+)|(v\d+alpha\d+))$`,
			},
		},
	})
	functional_test_helpers.TestRule(t, "use_import_alias_naming_custom_config_with_multiple_values", &use_import_alias_naming.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"allow-regex": `^[a-z][a-z0-9]*$`,
				"deny-regex":  `^((v\d+)|(v\d+alpha\d+))$`,
			},
		},
	})
}

func TestUseImportAliasNaming_CustomConfigWithOnlyDeny(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_import_alias_naming_custom_config_with_only_deny", &use_import_alias_naming.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"denyRegex": `^((v\d+)|(v\d+alpha\d+))$`,
			},
		},
	})
}

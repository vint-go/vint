package use_var_naming_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_var_naming"
)

func TestVarNaming(t *testing.T) {
	functional_test_helpers.TestRule(t, "var_naming", &use_var_naming.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{"ID"}, []any{"VM"}},
	})
	functional_test_helpers.TestRule(t, "var_naming_skip_initialism_name_checks_true", &use_var_naming.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			[]any{},
			[]any{},
			[]any{map[string]any{"skip-initialism-name-checks": true}},
		},
	})
	functional_test_helpers.TestRule(t, "var_naming_skip_initialism_name_checks_false", &use_var_naming.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			[]any{},
			[]any{},
			[]any{map[string]any{"skip-initialism-name-checks": false}},
		},
	})
	functional_test_helpers.TestRule(t, "var_naming_allowlist_blocklist_skip_initialism_name_checks", &use_var_naming.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			[]any{"ID"},
			[]any{"VM"},
			[]any{map[string]any{"skip-initialism-name-checks": true}},
		},
	})

	functional_test_helpers.TestRule(t, "var_naming_test", &use_var_naming.VarNamingRule{}, &lint.RuleConfig{})

	functional_test_helpers.TestRule(t, "var_naming_upper_case_const_false", &use_var_naming.VarNamingRule{}, &lint.RuleConfig{})
	functional_test_helpers.TestRule(t, "var_naming_upper_case_const_true", &use_var_naming.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{}, []any{}, []any{map[string]any{"upperCaseConst": true}}},
	})
	functional_test_helpers.TestRule(t, "var_naming_upper_case_const_true", &use_var_naming.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{}, []any{}, []any{map[string]any{"upper-case-const": true}}},
	})
}

func BenchmarkUpperCaseConstTrue(b *testing.B) {
	for b.Loop() {
		functional_test_helpers.TestRule(b, "var_naming_upper_case_const_true", &use_var_naming.VarNamingRule{}, &lint.RuleConfig{
			Arguments: lint.Arguments{[]any{}, []any{}, []any{map[string]any{"upperCaseConst": true}}},
		})
	}
}

package use_context_as_first_param_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_context_as_first_param"
)

func TestContextAsArgumentDefault(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_context_as_first_param_default", &use_context_as_first_param.ContextAsArgumentRule{})
}

func TestContextAsArgument(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_context_as_first_param", &use_context_as_first_param.ContextAsArgumentRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"allowTypesBefore": "AllowedBeforeType,AllowedBeforeStruct,*AllowedBeforePtrStruct,*testing.T",
			},
		},
	})
	functional_test_helpers.TestRule(t, "use_context_as_first_param", &use_context_as_first_param.ContextAsArgumentRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"allow-types-before": "AllowedBeforeType,AllowedBeforeStruct,*AllowedBeforePtrStruct,*testing.T",
			},
		},
	})
}

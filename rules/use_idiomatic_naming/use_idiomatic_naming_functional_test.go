package use_idiomatic_naming_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_idiomatic_naming"
)

func TestUseIdiomaticNaming(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_idiomatic_naming", &use_idiomatic_naming.UseIdiomaticNamingRule{})
}

func TestUseIdiomaticNamingCustomInitialisms(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_idiomatic_naming_custom_initialisms", &use_idiomatic_naming.UseIdiomaticNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"extraInitialisms": []any{"GRPC", "AMQP"},
			},
		},
	})
}

func TestUseIdiomaticNamingExcludeInitialisms(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_idiomatic_naming_exclude_initialisms", &use_idiomatic_naming.UseIdiomaticNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"excludeInitialisms": []any{"HTTP", "URL"},
			},
		},
	})
}

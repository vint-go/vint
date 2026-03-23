package use_map_style_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_map_style"
)

func TestUseMapStyle_any(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_map_style_any", &use_map_style.EnforceMapStyleRule{})
}

func TestUseMapStyle_make(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_map_style_make", &use_map_style.EnforceMapStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"make"},
	})
}

func TestUseMapStyle_literal(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_map_style_literal", &use_map_style.EnforceMapStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"literal"},
	})
}

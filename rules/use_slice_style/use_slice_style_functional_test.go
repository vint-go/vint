package use_slice_style_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_slice_style"
)

func TestUseSliceStyle_any(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_slice_style_any", &use_slice_style.EnforceSliceStyleRule{})
}

func TestUseSliceStyle_make(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_slice_style_make", &use_slice_style.EnforceSliceStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"make"},
	})
}

func TestUseSliceStyle_literal(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_slice_style_literal", &use_slice_style.EnforceSliceStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"literal"},
	})
}

func TestUseSliceStyle_nil(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_slice_style_nil", &use_slice_style.EnforceSliceStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"nil"},
	})
}

package use_inline_math_pow_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_inline_math_pow"
)

func TestUseInlineMathPow(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_inline_math_pow", &use_inline_math_pow.UseInlineMathPowRule{})
}

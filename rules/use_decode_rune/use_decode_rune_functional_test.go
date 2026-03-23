package use_decode_rune_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_decode_rune"
)

func TestUseDecodeRune(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_decode_rune", &use_decode_rune.UseDecodeRuneRule{})
}

package no_unused_constant_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unused_constant"
)

func TestNoUnusedConstant(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unused_constant", &no_unused_constant.NoUnusedConstantRule{})
}

package no_constant_result_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_constant_result"
)

func TestNoConstantResult(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_constant_result", &no_constant_result.NoConstantResultRule{})
}

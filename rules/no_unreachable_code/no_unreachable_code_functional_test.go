package no_unreachable_code_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unreachable_code"
)

func TestNoUnreachableCode(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unreachable_code", &no_unreachable_code.NoUnreachableCodeRule{})
}

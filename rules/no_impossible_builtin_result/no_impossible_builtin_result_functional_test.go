package no_impossible_builtin_result_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_impossible_builtin_result"
)

func TestNoImpossibleBuiltinResult(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_impossible_builtin_result", &no_impossible_builtin_result.NoImpossibleBuiltinResultRule{})
}

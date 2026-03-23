package no_always_true_len_check_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_always_true_len_check"
)

func TestNoAlwaysTrueLenCheck(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_always_true_len_check", &no_always_true_len_check.NoAlwaysTrueLenCheckRule{})
}

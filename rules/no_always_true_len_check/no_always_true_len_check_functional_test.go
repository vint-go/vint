package no_always_true_len_check_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_always_true_len_check"
)

func TestNoAlwaysTrueLenCheck(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_always_true_len_check", &no_always_true_len_check.NoAlwaysTrueLenCheckRule{})
}

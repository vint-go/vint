package no_std_method_signature_mismatch_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_std_method_signature_mismatch"
)

func TestNoStdMethodSignatureMismatch(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_std_method_signature_mismatch", &no_std_method_signature_mismatch.NoStdMethodSignatureMismatchRule{})
}

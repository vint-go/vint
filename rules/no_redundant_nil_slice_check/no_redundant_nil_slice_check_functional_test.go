package no_redundant_nil_slice_check_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_nil_slice_check"
)

func TestNoRedundantNilSliceCheck(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_nil_slice_check", &no_redundant_nil_slice_check.NoRedundantNilSliceCheckRule{})
}

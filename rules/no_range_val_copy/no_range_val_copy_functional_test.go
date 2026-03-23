package no_range_val_copy_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_range_val_copy"
)

func TestNoRangeValCopy(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_range_val_copy", &no_range_val_copy.NoRangeValCopyRule{})
}

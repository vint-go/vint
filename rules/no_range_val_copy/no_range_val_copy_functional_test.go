package no_range_val_copy_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_range_val_copy"
)

func TestNoRangeValCopy(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_range_val_copy", &no_range_val_copy.NoRangeValCopyRule{})
}

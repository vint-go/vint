package no_range_append_all_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_range_append_all"
)

func TestNoRangeAppendAll(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_range_append_all", &no_range_append_all.NoRangeAppendAllRule{})
}

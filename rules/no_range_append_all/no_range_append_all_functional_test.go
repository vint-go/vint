package no_range_append_all_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_range_append_all"
)

func TestNoRangeAppendAll(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_range_append_all", &no_range_append_all.NoRangeAppendAllRule{})
}

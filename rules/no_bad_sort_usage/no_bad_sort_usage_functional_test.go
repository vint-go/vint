package no_bad_sort_usage_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_bad_sort_usage"
)

func TestNoBadSortUsage(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_bad_sort_usage", &no_bad_sort_usage.NoBadSortUsageRule{})
}

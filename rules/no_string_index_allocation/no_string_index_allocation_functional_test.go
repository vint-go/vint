package no_string_index_allocation_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_string_index_allocation"
)

func TestNoStringIndexAllocation(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_string_index_allocation", &no_string_index_allocation.NoStringIndexAllocationRule{})
}

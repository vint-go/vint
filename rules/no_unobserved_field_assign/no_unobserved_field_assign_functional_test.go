package no_unobserved_field_assign_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unobserved_field_assign"
)

func TestNoUnobservedFieldAssign(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unobserved_field_assign", &no_unobserved_field_assign.NoUnobservedFieldAssignRule{})
}

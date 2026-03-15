package no_unreachable_type_case_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unreachable_type_case"
)

func TestNoUnreachableTypeCase(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unreachable_type_case", &no_unreachable_type_case.NoUnreachableTypeCaseRule{})
}

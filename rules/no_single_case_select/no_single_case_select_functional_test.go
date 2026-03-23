package no_single_case_select_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_single_case_select"
)

func TestNoSingleCaseSelect(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_single_case_select", &no_single_case_select.NoSingleCaseSelectRule{})
}

package no_single_case_select_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_single_case_select"
)

func TestNoSingleCaseSelect(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_single_case_select", &no_single_case_select.NoSingleCaseSelectRule{})
}

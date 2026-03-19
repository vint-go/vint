package no_confusing_results_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_confusing_results"
)

func TestNoConfusingResults(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_confusing_results", &no_confusing_results.ConfusingResultsRule{})
}

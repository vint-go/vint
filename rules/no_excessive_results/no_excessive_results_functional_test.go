package no_excessive_results_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_excessive_results"
)

func TestNoExcessiveResults(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_results", &no_excessive_results.NoExcessiveResultsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(3)},
	})
}

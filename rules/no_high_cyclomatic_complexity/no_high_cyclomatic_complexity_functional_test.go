package no_high_cyclomatic_complexity_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_high_cyclomatic_complexity"
)

func TestNoHighCyclomaticComplexity(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_high_cyclomatic_complexity", &no_high_cyclomatic_complexity.NoHighCyclomaticComplexityRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(5)},
	})
}

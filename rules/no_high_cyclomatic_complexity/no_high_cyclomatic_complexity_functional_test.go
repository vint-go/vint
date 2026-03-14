package no_high_cyclomatic_complexity_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_high_cyclomatic_complexity"
)

func TestNoHighCyclomaticComplexity(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_high_cyclomatic_complexity", &no_high_cyclomatic_complexity.NoHighCyclomaticComplexityRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(5)},
	})
}

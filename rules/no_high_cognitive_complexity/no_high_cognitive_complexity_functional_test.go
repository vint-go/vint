package no_high_cognitive_complexity_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_high_cognitive_complexity"
)

func TestNoHighCognitiveComplexityDefault(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_high_cognitive_complexity_default", &no_high_cognitive_complexity.NoHighCognitiveComplexityRule{}, &lint.RuleConfig{})
}

func TestNoHighCognitiveComplexity(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_high_cognitive_complexity", &no_high_cognitive_complexity.NoHighCognitiveComplexityRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(0)},
	})
}

package no_high_cognitive_complexity_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_high_cognitive_complexity"
)

func TestNoHighCognitiveComplexityDefault(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_high_cognitive_complexity_default", &no_high_cognitive_complexity.NoHighCognitiveComplexityRule{}, &lint.RuleConfig{})
}

func TestNoHighCognitiveComplexity(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_high_cognitive_complexity", &no_high_cognitive_complexity.NoHighCognitiveComplexityRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(0)},
	})
}

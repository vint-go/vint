package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestArgumentsLimitDefault(t *testing.T) {
	testRule(t, "argument_limit_default", &rule.ArgumentsLimitRule{}, &lint.RuleConfig{})
}

func TestArgumentsLimit(t *testing.T) {
	testRule(t, "argument_limit", &rule.ArgumentsLimitRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(3)},
	})
}

func BenchmarkArgumentsLimit(b *testing.B) {
	for b.Loop() {
		testRule(b, "argument_limit_default", &rule.ArgumentsLimitRule{}, &lint.RuleConfig{})
	}
}

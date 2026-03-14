package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestRangeValInClosure(t *testing.T) {
	testRule(t, "range_val_in_closure", &rule.RangeValInClosureRule{}, &lint.RuleConfig{})
}

func TestRangeValInClosureAfterGo1_22(t *testing.T) {
	testRule(t, "go1.22/range_val_in_closure", &rule.RangeValInClosureRule{}, &lint.RuleConfig{})
}

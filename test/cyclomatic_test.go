package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestCyclomaticDefault(t *testing.T) {
	testRule(t, "cyclomatic_default", &rule.CyclomaticRule{}, &lint.RuleConfig{})
}

func TestCyclomatic(t *testing.T) {
	testRule(t, "cyclomatic_default", &rule.CyclomaticRule{}, &lint.RuleConfig{})
	testRule(t, "cyclomatic", &rule.CyclomaticRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(1)},
	})
	testRule(t, "cyclomatic_2", &rule.CyclomaticRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(3)},
	})
}

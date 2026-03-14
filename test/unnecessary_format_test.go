package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestUnnecessaryFormat(t *testing.T) {
	testRule(t, "unnecessary_format", &rule.UnnecessaryFormatRule{}, &lint.RuleConfig{})
}

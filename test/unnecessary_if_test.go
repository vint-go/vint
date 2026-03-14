package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestUnnecessaryIf(t *testing.T) {
	testRule(t, "unnecessary_if", &rule.UnnecessaryIfRule{}, &lint.RuleConfig{})
}

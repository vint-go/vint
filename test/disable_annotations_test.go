package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestDisabledAnnotations(t *testing.T) {
	testRule(t, "disable_annotations", &rule.ExportedRule{}, &lint.RuleConfig{})
}

func TestModifiedAnnotations(t *testing.T) {
	testRule(t, "disable_annotations2", &rule.VarNamingRule{}, &lint.RuleConfig{})
}

func TestDisableNextLineAnnotations(t *testing.T) {
	testRule(t, "disable_annotations3", &rule.VarNamingRule{}, &lint.RuleConfig{})
}

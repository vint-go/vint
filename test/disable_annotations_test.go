package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/use_exported_comment"
	"github.com/strowk/vint/rules/use_var_naming"
)

func TestDisabledAnnotations(t *testing.T) {
	testRule(t, "disable_annotations", &use_exported_comment.ExportedRule{}, &lint.RuleConfig{})
}

func TestModifiedAnnotations(t *testing.T) {
	testRule(t, "disable_annotations2", &use_var_naming.VarNamingRule{}, &lint.RuleConfig{})
}

func TestDisableNextLineAnnotations(t *testing.T) {
	testRule(t, "disable_annotations3", &use_var_naming.VarNamingRule{}, &lint.RuleConfig{})
}

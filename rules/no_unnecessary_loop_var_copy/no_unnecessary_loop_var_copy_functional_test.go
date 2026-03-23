package no_unnecessary_loop_var_copy_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unnecessary_loop_var_copy"
)

func TestNoUnnecessaryLoopVarCopy(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_loop_var_copy", &no_unnecessary_loop_var_copy.NoUnnecessaryLoopVarCopyRule{})
}

func TestNoUnnecessaryLoopVarCopyCheckAlias(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_loop_var_copy_check_alias", &no_unnecessary_loop_var_copy.NoUnnecessaryLoopVarCopyRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"check-alias": true}},
	})
}

package no_unnecessary_loop_var_copy_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unnecessary_loop_var_copy"
)

func TestNoUnnecessaryLoopVarCopy(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_loop_var_copy", &no_unnecessary_loop_var_copy.NoUnnecessaryLoopVarCopyRule{})
}

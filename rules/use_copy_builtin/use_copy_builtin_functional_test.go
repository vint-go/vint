package use_copy_builtin_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_copy_builtin"
)

func TestUseCopyBuiltin(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_copy_builtin", &use_copy_builtin.UseCopyBuiltinRule{})
}

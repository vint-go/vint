package use_short_var_decl_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_short_var_decl"
)

func TestUseShortVarDecl(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_short_var_decl", &use_short_var_decl.UseShortVarDeclRule{})
}

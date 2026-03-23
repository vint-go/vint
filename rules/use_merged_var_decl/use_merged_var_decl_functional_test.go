package use_merged_var_decl_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_merged_var_decl"
)

func TestUseMergedVarDecl(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_merged_var_decl", &use_merged_var_decl.UseMergedVarDeclRule{})
}

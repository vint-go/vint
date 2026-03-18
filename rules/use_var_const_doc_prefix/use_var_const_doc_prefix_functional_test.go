package use_var_const_doc_prefix_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_var_const_doc_prefix"
)

func TestUseVarConstDocPrefix(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_var_const_doc_prefix", &use_var_const_doc_prefix.UseVarConstDocPrefixRule{})
}

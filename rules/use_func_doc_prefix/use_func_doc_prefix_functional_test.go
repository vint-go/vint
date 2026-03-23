package use_func_doc_prefix_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_func_doc_prefix"
)

func TestUseFuncDocPrefix(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_func_doc_prefix", &use_func_doc_prefix.UseFuncDocPrefixRule{})
}

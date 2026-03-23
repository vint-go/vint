package use_type_doc_prefix_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_type_doc_prefix"
)

func TestUseTypeDocPrefix(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_type_doc_prefix", &use_type_doc_prefix.UseTypeDocPrefixRule{})
}

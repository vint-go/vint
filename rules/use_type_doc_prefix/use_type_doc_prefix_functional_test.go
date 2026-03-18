package use_type_doc_prefix_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_type_doc_prefix"
)

func TestUseTypeDocPrefix(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_type_doc_prefix", &use_type_doc_prefix.UseTypeDocPrefixRule{})
}

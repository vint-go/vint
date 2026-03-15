package no_doc_comment_stub_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_doc_comment_stub"
)

func TestNoDocCommentStub(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_doc_comment_stub", &no_doc_comment_stub.NoDocCommentStubRule{})
}

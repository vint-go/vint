package no_commented_out_import_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_commented_out_import"
)

func TestNoCommentedOutImport(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_commented_out_import", &no_commented_out_import.NoCommentedOutImportRule{})
}

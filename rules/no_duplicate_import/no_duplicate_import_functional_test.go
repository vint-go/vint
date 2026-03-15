package no_duplicate_import_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_duplicate_import"
)

func TestNoDuplicateImport(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicate_import", &no_duplicate_import.NoDuplicateImportRule{})
}

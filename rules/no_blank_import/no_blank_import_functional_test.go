package no_blank_import_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_blank_import"
)

func TestNoBlankImport(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_blank_import", &no_blank_import.NoBlankImportRule{})
}

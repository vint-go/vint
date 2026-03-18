package no_dot_import_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_dot_import"
)

func TestNoDotImport(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_dot_import", &no_dot_import.NoDotImportRule{})
}

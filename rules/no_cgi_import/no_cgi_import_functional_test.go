package no_cgi_import_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_cgi_import"
)

func TestNoCgiImport(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_cgi_import", &no_cgi_import.NoCgiImportRule{})
}

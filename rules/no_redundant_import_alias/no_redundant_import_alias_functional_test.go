package no_redundant_import_alias_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_redundant_import_alias"
)

func TestNoRedundantImportAlias(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_import_alias", &no_redundant_import_alias.NoRedundantImportAliasRule{})
}

package no_redundant_import_alias_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_import_alias"
)

func TestNoRedundantImportAlias(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_import_alias", &no_redundant_import_alias.NoRedundantImportAliasRule{})
}

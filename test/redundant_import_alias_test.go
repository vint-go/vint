package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestRedundantImportAlias(t *testing.T) {
	testRule(t, "redundant_import_alias", &rule.RedundantImportAlias{})
}

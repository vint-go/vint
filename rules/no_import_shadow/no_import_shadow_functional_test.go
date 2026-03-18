package no_import_shadow_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_import_shadow"
)

func TestNoImportShadow(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_import_shadow", &no_import_shadow.NoImportShadowRule{})
}

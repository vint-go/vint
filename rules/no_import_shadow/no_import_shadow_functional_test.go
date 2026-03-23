package no_import_shadow_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_import_shadow"
)

func TestNoImportShadow(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_import_shadow", &no_import_shadow.NoImportShadowRule{})
}

package no_duplicated_imports_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_duplicated_imports"
)

func TestNoDuplicatedImports(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicated_imports", &no_duplicated_imports.DuplicatedImportsRule{})
}

package no_duplicated_imports_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_duplicated_imports"
)

func TestNoDuplicatedImports(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicated_imports", &no_duplicated_imports.DuplicatedImportsRule{})
}

package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestDuplicatedImports(t *testing.T) {
	testRule(t, "duplicated_imports", &rule.DuplicatedImportsRule{})
}

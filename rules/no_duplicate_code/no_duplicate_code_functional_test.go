package no_duplicate_code_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_duplicate_code"
)

func TestNoDuplicateCode(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicate_code", &no_duplicate_code.NoDuplicateCodeRule{})
}

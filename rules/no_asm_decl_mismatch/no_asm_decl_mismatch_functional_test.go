package no_asm_decl_mismatch_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_asm_decl_mismatch"
)

func TestNoAsmDeclMismatch(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_asm_decl_mismatch", &no_asm_decl_mismatch.NoAsmDeclMismatchRule{})
}

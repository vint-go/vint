package no_asm_decl_mismatch_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_asm_decl_mismatch"
)

func TestNoAsmDeclMismatch(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_asm_decl_mismatch", &no_asm_decl_mismatch.NoAsmDeclMismatchRule{})
}

package no_builtin_shadow_decl_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_builtin_shadow_decl"
)

func TestNoBuiltinShadowDecl(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_builtin_shadow_decl", &no_builtin_shadow_decl.NoBuiltinShadowDeclRule{})
}

package no_builtin_shadow_decl_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_builtin_shadow_decl"
)

func TestNoBuiltinShadowDecl(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_builtin_shadow_decl", &no_builtin_shadow_decl.NoBuiltinShadowDeclRule{})
}

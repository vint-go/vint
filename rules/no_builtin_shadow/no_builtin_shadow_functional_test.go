package no_builtin_shadow_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_builtin_shadow"
)

func TestNoBuiltinShadow(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_builtin_shadow", &no_builtin_shadow.NoBuiltinShadowRule{})
}

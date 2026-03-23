package use_else_if_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_else_if"
)

func TestUseElseIf(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_else_if", &use_else_if.UseElseIfRule{})
}

package use_simplified_bool_expr_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_simplified_bool_expr"
)

func TestUseSimplifiedBoolExpr(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_simplified_bool_expr", &use_simplified_bool_expr.UseSimplifiedBoolExprRule{})
}

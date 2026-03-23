package test_test

import (
	"testing"

	"github.com/vint-go/vint/rule"
)

func TestConstantLogicalExpr(t *testing.T) {
	testRule(t, "constant_logical_expr", &rule.ConstantLogicalExprRule{})
}

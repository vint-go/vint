package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestBoolLiteral(t *testing.T) {
	testRule(t, "bool_literal_in_expr", &rule.BoolLiteralRule{})
}

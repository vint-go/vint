package no_type_assert_else_misread_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_type_assert_else_misread"
)

func TestNoTypeAssertElseMisread(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_type_assert_else_misread", &no_type_assert_else_misread.NoTypeAssertElseMisreadRule{})
}

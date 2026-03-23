package no_impossible_interface_assert_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_impossible_interface_assert"
)

func TestNoImpossibleInterfaceAssert(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_impossible_interface_assert", &no_impossible_interface_assert.NoImpossibleInterfaceAssertRule{})
}

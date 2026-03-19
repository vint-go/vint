package no_nested_structs_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_nested_structs"
)

func TestNoNestedStructs(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_nested_structs", &no_nested_structs.NestedStructs{})
}

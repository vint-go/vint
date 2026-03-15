package no_magic_number_in_operation_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_magic_number_in_operation"
)

func TestNoMagicNumberInOperation(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_magic_number_in_operation", &no_magic_number_in_operation.NoMagicNumberInOperationRule{})
}

package no_unexported_return_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unexported_return"
)

func TestNoUnexportedReturn(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unexported_return_package_foo", &no_unexported_return.UnexportedReturnRule{})
	functional_test_helpers.TestRule(t, "no_unexported_return_package_foo_test", &no_unexported_return.UnexportedReturnRule{})
	functional_test_helpers.TestRule(t, "no_unexported_return_package_footest_test", &no_unexported_return.UnexportedReturnRule{})
	functional_test_helpers.TestRule(t, "no_unexported_return_package_main", &no_unexported_return.UnexportedReturnRule{})
	functional_test_helpers.TestRule(t, "no_unexported_return_package_main_test", &no_unexported_return.UnexportedReturnRule{})
	functional_test_helpers.TestRule(t, "no_unexported_return_package_maintest_test", &no_unexported_return.UnexportedReturnRule{})
}

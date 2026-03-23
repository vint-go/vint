package no_eval_order_dependency_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_eval_order_dependency"
)

func TestNoEvalOrderDependency(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_eval_order_dependency", &no_eval_order_dependency.NoEvalOrderDependencyRule{})
}

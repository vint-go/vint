package no_eval_order_dependency_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_eval_order_dependency"
)

func TestNoEvalOrderDependency(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_eval_order_dependency", &no_eval_order_dependency.NoEvalOrderDependencyRule{})
}

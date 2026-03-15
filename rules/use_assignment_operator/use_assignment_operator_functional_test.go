package use_assignment_operator_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_assignment_operator"
)

func TestUseAssignmentOperator(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_assignment_operator", &use_assignment_operator.UseAssignmentOperatorRule{})
}

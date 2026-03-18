package no_duplicate_build_constraint_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_duplicate_build_constraint"
)

func TestNoDuplicateBuildConstraint(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicate_build_constraint", &no_duplicate_build_constraint.NoDuplicateBuildConstraintRule{})
}

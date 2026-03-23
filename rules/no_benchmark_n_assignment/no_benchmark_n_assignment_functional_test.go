package no_benchmark_n_assignment_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_benchmark_n_assignment"
)

func TestNoBenchmarkNAssignment(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_benchmark_n_assignment", &no_benchmark_n_assignment.NoBenchmarkNAssignmentRule{})
}

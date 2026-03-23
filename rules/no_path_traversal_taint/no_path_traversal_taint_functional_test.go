package no_path_traversal_taint_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_path_traversal_taint"
)

func TestNoPathTraversalTaint(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_path_traversal_taint", &no_path_traversal_taint.NoPathTraversalTaintRule{})
}

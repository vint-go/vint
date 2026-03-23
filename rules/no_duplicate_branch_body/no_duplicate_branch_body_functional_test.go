package no_duplicate_branch_body_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_duplicate_branch_body"
)

func TestNoDuplicateBranchBody(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_duplicate_branch_body", &no_duplicate_branch_body.NoDuplicateBranchBodyRule{})
}

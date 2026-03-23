package no_redundant_sprint_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_redundant_sprint"
)

func TestNoRedundantSprint(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_redundant_sprint", &no_redundant_sprint.NoRedundantSprintRule{})
}

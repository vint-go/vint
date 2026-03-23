package no_space_in_directive_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_space_in_directive"
)

func TestNoSpaceInDirective(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_space_in_directive", &no_space_in_directive.NoSpaceInDirectiveRule{})
}

package no_unrecognized_directive_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unrecognized_directive"
)

func TestNoUnrecognizedDirective(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unrecognized_directive", &no_unrecognized_directive.NoUnrecognizedDirectiveRule{})
}

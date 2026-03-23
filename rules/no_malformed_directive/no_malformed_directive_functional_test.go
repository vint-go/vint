package no_malformed_directive_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_malformed_directive"
)

func TestNoMalformedDirective(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_malformed_directive", &no_malformed_directive.NoMalformedDirectiveRule{})
}

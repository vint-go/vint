package no_malformed_directive_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_malformed_directive"
)

func TestNoMalformedDirective(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_malformed_directive", &no_malformed_directive.NoMalformedDirectiveRule{})
}

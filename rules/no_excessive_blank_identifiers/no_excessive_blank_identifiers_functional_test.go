package no_excessive_blank_identifiers_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_excessive_blank_identifiers"
)

func TestNoExcessiveBlankIdentifiers(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_excessive_blank_identifiers", &no_excessive_blank_identifiers.NoExcessiveBlankIdentifiersRule{})
}

package no_unnecessary_blank_identifier_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unnecessary_blank_identifier"
)

func TestNoUnnecessaryBlankIdentifier(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_blank_identifier", &no_unnecessary_blank_identifier.NoUnnecessaryBlankIdentifierRule{})
}

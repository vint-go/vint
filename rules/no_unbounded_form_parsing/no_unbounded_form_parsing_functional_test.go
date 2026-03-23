package no_unbounded_form_parsing_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unbounded_form_parsing"
)

func TestNoUnboundedFormParsing(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unbounded_form_parsing", &no_unbounded_form_parsing.NoUnboundedFormParsingRule{})
}

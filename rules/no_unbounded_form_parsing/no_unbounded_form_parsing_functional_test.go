package no_unbounded_form_parsing_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unbounded_form_parsing"
)

func TestNoUnboundedFormParsing(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unbounded_form_parsing", &no_unbounded_form_parsing.NoUnboundedFormParsingRule{})
}

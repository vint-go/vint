package no_invalid_template_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_invalid_template"
)

func TestNoInvalidTemplate(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_invalid_template", &no_invalid_template.NoInvalidTemplateRule{})
}

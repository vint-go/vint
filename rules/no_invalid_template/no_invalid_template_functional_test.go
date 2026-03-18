package no_invalid_template_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_invalid_template"
)

func TestNoInvalidTemplate(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_invalid_template", &no_invalid_template.NoInvalidTemplateRule{})
}

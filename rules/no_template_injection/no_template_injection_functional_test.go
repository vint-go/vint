package no_template_injection_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_template_injection"
)

func TestNoTemplateInjection(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_template_injection", &no_template_injection.NoTemplateInjectionRule{})
}

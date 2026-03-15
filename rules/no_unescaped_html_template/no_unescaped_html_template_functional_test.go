package no_unescaped_html_template_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unescaped_html_template"
)

func TestNoUnescapedHtmlTemplate(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unescaped_html_template", &no_unescaped_html_template.NoUnescapedHtmlTemplateRule{})
}

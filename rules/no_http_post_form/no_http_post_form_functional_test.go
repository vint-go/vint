package no_http_post_form_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_http_post_form"
)

func TestNoHttpPostForm(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_http_post_form", &no_http_post_form.NoHttpPostFormRule{})
}

package no_unsecure_url_scheme_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unsecure_url_scheme"
)

func TestNoUnsecureUrlScheme(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unsecure_url_scheme", &no_unsecure_url_scheme.UnsecureURLSchemeRule{})
}

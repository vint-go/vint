package no_http_response_misuse_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_http_response_misuse"
)

func TestNoHttpResponseMisuse(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_http_response_misuse", &no_http_response_misuse.NoHttpResponseMisuseRule{})
}

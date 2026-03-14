package no_http_request_smuggling_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_http_request_smuggling"
)

func TestNoHttpRequestSmuggling(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_http_request_smuggling", &no_http_request_smuggling.NoHttpRequestSmugglingRule{})
}

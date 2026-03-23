package no_http_request_smuggling_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_http_request_smuggling"
)

func TestNoHttpRequestSmuggling(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_http_request_smuggling", &no_http_request_smuggling.NoHttpRequestSmugglingRule{})
}

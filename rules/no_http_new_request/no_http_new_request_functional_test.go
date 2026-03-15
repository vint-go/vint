package no_http_new_request_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_http_new_request"
)

func TestNoHttpNewRequest(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_http_new_request", &no_http_new_request.NoHttpNewRequestRule{})
}

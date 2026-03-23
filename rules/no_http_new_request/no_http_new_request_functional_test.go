package no_http_new_request_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_http_new_request"
)

func TestNoHttpNewRequest(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_http_new_request", &no_http_new_request.NoHttpNewRequestRule{})
}

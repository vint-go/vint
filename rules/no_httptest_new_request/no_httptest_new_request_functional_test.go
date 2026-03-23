package no_httptest_new_request_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_httptest_new_request"
)

func TestNoHttptestNewRequest(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_httptest_new_request", &no_httptest_new_request.NoHttptestNewRequestRule{})
}

package no_httptest_new_request_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_httptest_new_request"
)

func TestNoHttptestNewRequest(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_httptest_new_request", &no_httptest_new_request.NoHttptestNewRequestRule{})
}

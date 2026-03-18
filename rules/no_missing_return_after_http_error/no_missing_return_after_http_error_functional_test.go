package no_missing_return_after_http_error_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_missing_return_after_http_error"
)

func TestNoMissingReturnAfterHttpError(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_missing_return_after_http_error", &no_missing_return_after_http_error.NoMissingReturnAfterHttpErrorRule{})
}

package no_off_by_one_error_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_off_by_one_error"
)

func TestNoOffByOneError(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_off_by_one_error", &no_off_by_one_error.NoOffByOneErrorRule{})
}

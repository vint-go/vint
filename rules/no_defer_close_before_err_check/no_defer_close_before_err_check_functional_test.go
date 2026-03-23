package no_defer_close_before_err_check_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_defer_close_before_err_check"
)

func TestNoDeferCloseBeforeErrCheck(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_defer_close_before_err_check", &no_defer_close_before_err_check.NoDeferCloseBeforeErrCheckRule{})
}

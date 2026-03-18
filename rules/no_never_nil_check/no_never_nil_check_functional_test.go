package no_never_nil_check_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_never_nil_check"
)

func TestNoNeverNilCheck(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_never_nil_check", &no_never_nil_check.NoNeverNilCheckRule{})
}

package no_timer_reset_retval_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_timer_reset_retval"
)

func TestNoTimerResetRetval(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_timer_reset_retval", &no_timer_reset_retval.NoTimerResetRetvalRule{})
}

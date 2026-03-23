package no_time_tick_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_time_tick"
)

func TestNoTimeTick(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_time_tick", &no_time_tick.NoTimeTickRule{})
}

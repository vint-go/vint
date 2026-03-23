package use_time_sleep_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_time_sleep"
)

func TestUseTimeSleep(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_time_sleep", &use_time_sleep.UseTimeSleepRule{})
}

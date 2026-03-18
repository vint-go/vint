package no_suspicious_time_sleep_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_suspicious_time_sleep"
)

func TestNoSuspiciousTimeSleep(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_suspicious_time_sleep", &no_suspicious_time_sleep.NoSuspiciousTimeSleepRule{})
}

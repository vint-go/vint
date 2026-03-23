package no_suspicious_time_sleep_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_suspicious_time_sleep"
)

func TestNoSuspiciousTimeSleep(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_suspicious_time_sleep", &no_suspicious_time_sleep.NoSuspiciousTimeSleepRule{})
}

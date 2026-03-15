package no_defer_time_misuse_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_defer_time_misuse"
)

func TestNoDeferTimeMisuse(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_defer_time_misuse", &no_defer_time_misuse.NoDeferTimeMisuseRule{})
}

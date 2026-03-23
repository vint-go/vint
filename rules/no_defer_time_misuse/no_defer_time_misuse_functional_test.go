package no_defer_time_misuse_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_defer_time_misuse"
)

func TestNoDeferTimeMisuse(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_defer_time_misuse", &no_defer_time_misuse.NoDeferTimeMisuseRule{})
}

package no_deprecated_usage_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_deprecated_usage"
)

func TestNoDeprecatedUsage(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_deprecated_usage", &no_deprecated_usage.NoDeprecatedUsageRule{})
}

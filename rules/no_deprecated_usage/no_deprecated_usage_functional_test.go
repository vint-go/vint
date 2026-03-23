package no_deprecated_usage_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_deprecated_usage"
)

func TestNoDeprecatedUsage(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_deprecated_usage", &no_deprecated_usage.NoDeprecatedUsageRule{})
}

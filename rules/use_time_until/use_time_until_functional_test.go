package use_time_until_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_time_until"
)

func TestUseTimeUntil(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_time_until", &use_time_until.UseTimeUntilRule{})
}

package use_time_until_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_time_until"
)

func TestUseTimeUntil(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_time_until", &use_time_until.UseTimeUntilRule{})
}

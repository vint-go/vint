package use_time_since_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_time_since"
)

func TestUseTimeSince(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_time_since", &use_time_since.UseTimeSinceRule{})
}

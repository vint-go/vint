package use_percent_q_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_percent_q"
)

func TestUsePercentQ(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_percent_q", &use_percent_q.UsePercentQRule{})
}

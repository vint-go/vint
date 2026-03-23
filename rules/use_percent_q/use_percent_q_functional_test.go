package use_percent_q_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_percent_q"
)

func TestUsePercentQ(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_percent_q", &use_percent_q.UsePercentQRule{})
}

package use_consistent_receiver_name_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_consistent_receiver_name"
)

func TestUseConsistentReceiverName(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_consistent_receiver_name", &use_consistent_receiver_name.UseConsistentReceiverNameRule{})
}

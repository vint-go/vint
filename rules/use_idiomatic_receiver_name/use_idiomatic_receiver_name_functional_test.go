package use_idiomatic_receiver_name_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_idiomatic_receiver_name"
)

func TestUseIdiomaticReceiverName(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_idiomatic_receiver_name", &use_idiomatic_receiver_name.UseIdiomaticReceiverNameRule{})
}

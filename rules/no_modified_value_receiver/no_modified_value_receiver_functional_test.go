package no_modified_value_receiver_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_modified_value_receiver"
)

func TestNoModifiedValueReceiver(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_modified_value_receiver", &no_modified_value_receiver.ModifiesValRecRule{})
}

package no_lost_cancel_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_lost_cancel"
)

func TestNoLostCancel(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_lost_cancel", &no_lost_cancel.NoLostCancelRule{})
}

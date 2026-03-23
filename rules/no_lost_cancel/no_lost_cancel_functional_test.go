package no_lost_cancel_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_lost_cancel"
)

func TestNoLostCancel(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_lost_cancel", &no_lost_cancel.NoLostCancelRule{})
}

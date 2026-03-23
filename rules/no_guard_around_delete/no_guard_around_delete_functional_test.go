package no_guard_around_delete_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_guard_around_delete"
)

func TestNoGuardAroundDelete(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_guard_around_delete", &no_guard_around_delete.NoGuardAroundDeleteRule{})
}

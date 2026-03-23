package no_copied_lock_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_copied_lock"
)

func TestNoCopiedLock(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_copied_lock", &no_copied_lock.NoCopiedLockRule{})
}

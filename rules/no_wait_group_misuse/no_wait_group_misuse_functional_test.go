package no_wait_group_misuse_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_wait_group_misuse"
)

func TestNoWaitGroupMisuse(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_wait_group_misuse", &no_wait_group_misuse.NoWaitGroupMisuseRule{})
}

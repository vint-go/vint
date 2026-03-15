package no_inline_sync_once_func_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_inline_sync_once_func"
)

func TestNoInlineSyncOnceFunc(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_inline_sync_once_func", &no_inline_sync_once_func.NoInlineSyncOnceFuncRule{})
}

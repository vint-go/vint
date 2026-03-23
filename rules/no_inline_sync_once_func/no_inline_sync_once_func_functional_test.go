package no_inline_sync_once_func_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_inline_sync_once_func"
)

func TestNoInlineSyncOnceFunc(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_inline_sync_once_func", &no_inline_sync_once_func.NoInlineSyncOnceFuncRule{})
}

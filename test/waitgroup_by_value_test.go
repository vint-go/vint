package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestWaitGroupByValue(t *testing.T) {
	testRule(t, "waitgroup_by_value", &rule.WaitGroupByValueRule{})
}

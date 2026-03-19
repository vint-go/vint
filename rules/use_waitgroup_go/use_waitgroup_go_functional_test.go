package use_waitgroup_go_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_waitgroup_go"
)

func TestUseWaitgroupGo(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_waitgroup_go", &use_waitgroup_go.UseWaitGroupGoRule{})
	functional_test_helpers.TestRule(t, "go1.25/use_waitgroup_go", &use_waitgroup_go.UseWaitGroupGoRule{})
}

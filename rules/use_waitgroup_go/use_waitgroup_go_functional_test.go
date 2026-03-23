package use_waitgroup_go_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_waitgroup_go"
)

func TestUseWaitgroupGo(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_waitgroup_go", &use_waitgroup_go.UseWaitGroupGoRule{})
	functional_test_helpers.TestRule(t, "go1.25/use_waitgroup_go", &use_waitgroup_go.UseWaitGroupGoRule{})
}

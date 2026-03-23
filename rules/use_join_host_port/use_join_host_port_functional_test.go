package use_join_host_port_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_join_host_port"
)

func TestUseJoinHostPort(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_join_host_port", &use_join_host_port.UseJoinHostPortRule{})
}

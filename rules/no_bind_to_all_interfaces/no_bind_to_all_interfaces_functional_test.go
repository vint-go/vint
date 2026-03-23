package no_bind_to_all_interfaces_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_bind_to_all_interfaces"
)

func TestNoBindToAllInterfaces(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_bind_to_all_interfaces", &no_bind_to_all_interfaces.NoBindToAllInterfacesRule{})
}

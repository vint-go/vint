package no_permissive_os_create_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_permissive_os_create"
)

func TestNoPermissiveOsCreate(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_permissive_os_create", &no_permissive_os_create.NoPermissiveOsCreateRule{})
}

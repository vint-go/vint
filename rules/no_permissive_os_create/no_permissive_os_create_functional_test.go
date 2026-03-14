package no_permissive_os_create_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_permissive_os_create"
)

func TestNoPermissiveOsCreate(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_permissive_os_create", &no_permissive_os_create.NoPermissiveOsCreateRule{})
}

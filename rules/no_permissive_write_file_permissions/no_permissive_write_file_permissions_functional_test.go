package no_permissive_write_file_permissions_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_permissive_write_file_permissions"
)

func TestNoPermissiveWriteFilePermissions(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_permissive_write_file_permissions", &no_permissive_write_file_permissions.NoPermissiveWriteFilePermissionsRule{})
}

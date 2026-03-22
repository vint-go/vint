package no_permissive_directory_permissions_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_permissive_directory_permissions"
)

func TestNoPermissiveDirectoryPermissions(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_permissive_directory_permissions", &no_permissive_directory_permissions.NoPermissiveDirectoryPermissionsRule{})
}

func TestNoPermissiveDirectoryPermissionsCustom(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_permissive_directory_permissions_custom", &no_permissive_directory_permissions.NoPermissiveDirectoryPermissionsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"maxPermission": int64(0o700),
		}},
	})
}

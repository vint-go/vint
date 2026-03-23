package no_permissive_directory_permissions_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_permissive_directory_permissions"
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

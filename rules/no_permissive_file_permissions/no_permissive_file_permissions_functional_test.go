package no_permissive_file_permissions_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_permissive_file_permissions"
)

func TestNoPermissiveFilePermissions(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_permissive_file_permissions", &no_permissive_file_permissions.NoPermissiveFilePermissionsRule{})
}

func TestNoPermissiveFilePermissionsCustom(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_permissive_file_permissions_custom", &no_permissive_file_permissions.NoPermissiveFilePermissionsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"maxPermission": int64(0o644),
		}},
	})
}

package no_file_scoped_denied_import_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_file_scoped_denied_import"
)

func TestNoFileScopedDeniedImport(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_file_scoped_denied_import", &no_file_scoped_denied_import.NoFileScopedDeniedImportRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"rules": map[string]any{
				"main": map[string]any{
					"files": []any{"!$test"},
					"deny": []any{
						map[string]any{
							"pkg":  "github.com/stretchr/testify",
							"desc": "testify should only be used in test files",
						},
					},
				},
			},
		}},
	})
}

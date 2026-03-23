package no_unallowed_import_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_unallowed_import"
)

func TestNoUnallowedImport(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unallowed_import", &no_unallowed_import.NoUnallowedImportRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"rules": map[string]any{
				"main": map[string]any{
					"list-mode": "strict",
					"files": []any{
						"$all",
					},
					"allow": []any{
						"$gostd",
						"github.com/myorg",
					},
				},
			},
		}},
	})
}

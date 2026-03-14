package no_denied_import_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_denied_import"
)

func TestNoDeniedImport(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_denied_import", &no_denied_import.NoDeniedImportRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"main": map[string]any{
				"list-mode": "lax",
				"deny": []any{
					map[string]any{
						"pkg":  "github.com/sirupsen/logrus",
						"desc": "Use log/slog instead of logrus",
					},
					map[string]any{
						"pkg":  "github.com/pkg/errors",
						"desc": "Use fmt.Errorf with %w verb for error wrapping",
					},
					map[string]any{
						"pkg":  "io/ioutil",
						"desc": "Use os and io packages directly instead",
					},
				},
			},
		}},
	})
}

package no_dot_import_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_dot_import"
)

func TestNoDotImport(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_dot_import", &no_dot_import.NoDotImportRule{})
}

func TestNoDotImportWithAllowedPackages(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_dot_import_allowed", &no_dot_import.NoDotImportRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"allowedPackages": []any{"errors", "context", "github.com/BurntSushi/toml"},
		}},
	})
}

func TestNoDotImportWithAllowedPackagesKebabCase(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_dot_import_allowed", &no_dot_import.NoDotImportRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"allowed-packages": []any{"errors", "context", "github.com/BurntSushi/toml"},
		}},
	})
}

package use_package_naming_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_package_naming"
)

func TestUsePackageNaming_conventionName(t *testing.T) {
	functional_test_helpers.TestRule(t, "package_naming_mixed_caps", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{})
	functional_test_helpers.TestRule(t, "package_naming_mixed_caps", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-convention-name-check": false},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_mixed_caps_skip", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-convention-name-check": true},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_underscore", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{})
	functional_test_helpers.TestRule(t, "package_naming_underscore", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-convention-name-check": false},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_underscore_skip", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-convention-name-check": true},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_convention_name_regex", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"skip-convention-name-check":  false,
				"convention-name-check-regex": "^[a-z][a-zA-Z0-9]*$", // allow camel case
			},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_convention_name_regex", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"convention-name-check-regex": "^[a-z][a-zA-Z0-9]*$", // allow camel case
			},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_convention_name_regex_test", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"convention-name-check-regex": "^[a-z][a-zA-Z0-9]*$", // allow camel case
			},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_convention_name_regex_skip", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"convention-name-check-regex": "^[a-z][a-zA-Z0-9]*$", // allow camel case
			},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_convention_name_regex_test_skip", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"convention-name-check-regex": "^[a-z][a-zA-Z0-9]*$", // allow camel case
			},
		},
	})
}

func TestUsePackageNaming_topLevel(t *testing.T) {
	functional_test_helpers.TestRule(t, "package_naming_top_level_pkg", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{})
	functional_test_helpers.TestRule(t, "package_naming_top_level_pkg", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-top-level-check": false},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_top_level_pkg_skip", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-top-level-check": true},
		},
	})
}

func TestUsePackageNaming_badNames(t *testing.T) {
	functional_test_helpers.TestRule(t, "package_naming_bad_default", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{})
	functional_test_helpers.TestRule(t, "package_naming_bad_default", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-default-bad-name-check": false},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_bad_default_skip", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-default-bad-name-check": true},
		},
	})

	functional_test_helpers.TestRule(t, "package_naming_bad_extra", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"check-extra-bad-name": true},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_bad_extra_skip", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{})
	functional_test_helpers.TestRule(t, "package_naming_bad_extra_skip", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"check-extra-bad-name": false},
		},
	})

	functional_test_helpers.TestRule(t, "package_naming_bad_user_defined", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"user-defined-bad-names": []any{"data"}},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_bad_user_defined_skip", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{})
	functional_test_helpers.TestRule(t, "package_naming_bad_user_defined_skip", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"user-defined-bad-names": []any{}},
		},
	})
}

func TestUsePackageNaming_stdLibConflict(t *testing.T) {
	functional_test_helpers.TestRule(t, "package_naming_std_common_conflict", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{})
	functional_test_helpers.TestRule(t, "package_naming_std_common_conflict", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-collision-with-common-std": false},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_std_common_conflict_skip", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-collision-with-common-std": true},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_std_all_conflict", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"check-collision-with-all-std": true},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_std_all_common", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"check-collision-with-all-std": true},
		},
	})
	functional_test_helpers.TestRule(t, "package_naming_std_all_conflict_skip", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{})
	functional_test_helpers.TestRule(t, "package_naming_std_all_conflict_skip", &use_package_naming.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"check-collision-with-all-std": false},
		},
	})
}

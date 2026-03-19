package use_repeated_arg_type_style_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_repeated_arg_type_style"
)

func TestUseRepeatedArgTypeStyleDefault(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_default", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{})
}

func TestUseRepeatedArgTypeStyleShort(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_short_args", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"short"},
	})
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_short_return", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"short"},
	})

	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_short_args", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"funcArgStyle": `short`,
			},
		},
	})
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_short_args", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"func-arg-style": `short`,
			},
		},
	})
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_short_return", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"funcRetValStyle": `short`,
			},
		},
	})
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_short_return", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"func-ret-val-style": `short`,
			},
		},
	})
}

func TestUseRepeatedArgTypeStyleFull(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_full_args", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"full"},
	})
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_full_return", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"full"},
	})

	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_full_args", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"funcArgStyle": `full`,
			},
		},
	})
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_full_return", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"funcRetValStyle": `full`,
			},
		},
	})
}

func TestUseRepeatedArgTypeStyleMixed(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_full_args", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"funcArgStyle": `full`,
			},
		},
	})
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_full_args", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"funcArgStyle":    `full`,
				"funcRetValStyle": `any`,
			},
		},
	})
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_full_args", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"funcArgStyle":    `full`,
				"funcRetValStyle": `short`,
			},
		},
	})

	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_full_return", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"funcRetValStyle": `full`,
			},
		},
	})
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_full_return", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"funcArgStyle":    `any`,
				"funcRetValStyle": `full`,
			},
		},
	})
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_full_return", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"funcArgStyle":    `short`,
				"funcRetValStyle": `full`,
			},
		},
	})

	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_mixed_full_short", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"funcArgStyle":    `full`,
				"funcRetValStyle": `short`,
			},
		},
	})
	functional_test_helpers.TestRule(t, "use_repeated_arg_type_style_mixed_short_full", &use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"funcArgStyle":    `short`,
				"funcRetValStyle": `full`,
			},
		},
	})
}

package use_string_format_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_string_format"
)

func TestUseStringFormat(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_string_format", &use_string_format.StringFormatRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			[]any{
				"stringFormatMethod1", // The first argument is checked by default
				"/^[A-Z]/",
				"must start with a capital letter",
			},

			[]any{
				"stringFormatMethod2[2].d",
				"/[^\\.]$/",
			}, // Must not end with a period
			[]any{
				"s.Method3[2]",
				"!/^[Tt][Hh]/",
				"must not start with 'th'",
			},
			[]any{
				"s.Method4", // same as before, but called from a struct
				"!/^[Ot][Tt]/",
				"must not start with 'ot'",
			},
		},
	})
}

func TestUseStringFormatDuplicatedStrings(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_string_format_issue_1063", &use_string_format.StringFormatRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{
			"fmt.Errorf[0],errors.New[0]",
			"/^([^A-Z]|$)/",
			"must not start with a capital letter",
		}},
	})
}

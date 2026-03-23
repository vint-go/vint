package no_error_strings_test

import (
	"errors"
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_error_strings"
)

func TestErrorStringsWithCustomFunctions(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_error_strings_with_custom_functions", &no_error_strings.ErrorStringsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"pkgErrors.Wrap"},
	})
}

func TestErrorStringsIssue1243(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_error_strings_issue_1243", &no_error_strings.ErrorStringsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"errors.Wrap"},
	})
}

func TestErrorStringsRule_Configure(t *testing.T) {
	tests := []struct {
		name      string
		arguments lint.Arguments
		wantErr   error
	}{
		{
			name:      "Default configuration",
			arguments: lint.Arguments{},
		},
		{
			name:      "Valid custom functions",
			arguments: lint.Arguments{"mypkg.MyErrorFunc", "errors.New"},
		},
		{
			name:      "Argument not a string",
			arguments: lint.Arguments{123},
		},
		{
			name:      "Invalid package",
			arguments: lint.Arguments{".MyErrorFunc"},
			wantErr:   errors.New("found invalid custom function: .MyErrorFunc"),
		},
		{
			name:      "Invalid function",
			arguments: lint.Arguments{"errors."},
			//revive:disable-next-line:lint/style/noErrorStrings
			wantErr: errors.New("found invalid custom function: errors."),
		},
		{
			name:      "Invalid custom function",
			arguments: lint.Arguments{"invalidFunction"},
			wantErr:   errors.New("found invalid custom function: invalidFunction"),
		},
		{
			name:      "Mixed valid and invalid custom functions",
			arguments: lint.Arguments{"mypkg.MyErrorFunc", "invalidFunction", "invalidFunction2"},
			wantErr:   errors.New("found invalid custom function: invalidFunction,invalidFunction2"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r no_error_strings.ErrorStringsRule

			err := r.Configure(tt.arguments)

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Configure() unexpected non-nil error %q", err)
				}
				return
			}
			if err == nil || err.Error() != tt.wantErr.Error() {
				t.Errorf("Configure() unexpected error: got %q, want %q", err, tt.wantErr)
			}
		})
	}
}

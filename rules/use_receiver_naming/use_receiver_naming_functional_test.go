package use_receiver_naming_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_receiver_naming"
)

func TestReceiverNamingTypeParams(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_receiver_naming_issue_669", &use_receiver_naming.ReceiverNamingRule{})
}

func TestReceiverNamingMaxLength(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_receiver_naming_issue_1040", &use_receiver_naming.ReceiverNamingRule{},
		&lint.RuleConfig{
			Arguments: lint.Arguments{
				map[string]any{"maxLength": int64(2)},
			},
		},
	)
	functional_test_helpers.TestRule(t, "use_receiver_naming_issue_1040", &use_receiver_naming.ReceiverNamingRule{},
		&lint.RuleConfig{
			Arguments: lint.Arguments{
				map[string]any{"max-length": int64(2)},
			},
		},
	)
}

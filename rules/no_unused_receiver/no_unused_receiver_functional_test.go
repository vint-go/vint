package no_unused_receiver_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unused_receiver"
)

func TestNoUnusedReceiver(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unused_receiver", &no_unused_receiver.UnusedReceiverRule{})
	functional_test_helpers.TestRule(t, "no_unused_receiver", &no_unused_receiver.UnusedReceiverRule{}, &lint.RuleConfig{Arguments: lint.Arguments{}})
	functional_test_helpers.TestRule(t, "no_unused_receiver", &no_unused_receiver.UnusedReceiverRule{}, &lint.RuleConfig{Arguments: lint.Arguments{
		map[string]any{"a": "^xxx"},
	}})
	functional_test_helpers.TestRule(t, "no_unused_receiver_custom_regex", &no_unused_receiver.UnusedReceiverRule{}, &lint.RuleConfig{Arguments: lint.Arguments{
		map[string]any{"allowRegex": "^xxx"},
	}})
	functional_test_helpers.TestRule(t, "no_unused_receiver_custom_regex", &no_unused_receiver.UnusedReceiverRule{}, &lint.RuleConfig{Arguments: lint.Arguments{
		map[string]any{"allow-regex": "^xxx"},
	}})
}

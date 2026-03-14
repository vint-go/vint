package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestIndentErrorFlow(t *testing.T) {
	testRule(t, "indent_error_flow", &rule.IndentErrorFlowRule{})
	testRule(t, "indent_error_flow_scope", &rule.IndentErrorFlowRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserveScope"}})
	testRule(t, "indent_error_flow_scope", &rule.IndentErrorFlowRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserve-scope"}})
}

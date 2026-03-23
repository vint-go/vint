package use_indent_error_flow_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/use_indent_error_flow"
)

func TestUseIndentErrorFlow(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_indent_error_flow", &use_indent_error_flow.IndentErrorFlowRule{})
	functional_test_helpers.TestRule(t, "use_indent_error_flow_scope", &use_indent_error_flow.IndentErrorFlowRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserveScope"}})
	functional_test_helpers.TestRule(t, "use_indent_error_flow_scope", &use_indent_error_flow.IndentErrorFlowRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserve-scope"}})
}

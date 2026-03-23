package no_variable_command_execution_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_variable_command_execution"
)

func TestNoVariableCommandExecution(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_variable_command_execution", &no_variable_command_execution.NoVariableCommandExecutionRule{})
}

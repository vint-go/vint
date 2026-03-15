package no_exec_command_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_exec_command"
)

func TestNoExecCommand(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_exec_command", &no_exec_command.NoExecCommandRule{})
}

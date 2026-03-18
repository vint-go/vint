package fixtures

import (
	"context"
	"os/exec"
)

func badExecCommandWithSpace() {
	cmd := exec.Command("echo hello") // MATCH /first argument to exec.Command looks like a shell command, but exec.Command expects the executable path without arguments/
	_ = cmd
}

func badExecCommandContextWithSpace() {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "echo hello") // MATCH /first argument to exec.Command looks like a shell command, but exec.Command expects the executable path without arguments/
	_ = cmd
}

func goodExecCommand() {
	cmd := exec.Command("echo", "hello")
	_ = cmd
}

func goodExecCommandContext() {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "echo", "hello")
	_ = cmd
}

func goodExecCommandSingleArg() {
	cmd := exec.Command("ls")
	_ = cmd
}

func goodExecCommandVariable() {
	name := "echo hello"
	cmd := exec.Command(name)
	_ = cmd
	_ = name
}

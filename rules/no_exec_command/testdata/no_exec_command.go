package fixtures

import (
	"context"
	"fmt"
	"os/exec"
)

func badExecCommand() {
	cmd := exec.Command("ls", "-la") // MATCH /exec.Command does not accept a context; use exec.CommandContext instead/
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))
}

func goodExecCommandContext() {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "ls", "-la")
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))
}

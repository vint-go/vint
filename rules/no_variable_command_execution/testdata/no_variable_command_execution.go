package fixtures

import (
	"context"
	"os/exec"
	"syscall"
)

// Invalid: Command name from variable
func runCommand(userInput string) error {
	cmd := exec.Command(userInput) // MATCH /potential command injection: variable argument passed to exec.Command/
	return cmd.Run()
}

// Invalid: Argument from variable - potential command injection
func processFile(filename string) ([]byte, error) {
	return exec.Command("cat", filename).Output() // MATCH /potential command injection: variable argument passed to exec.Command/
}

// Invalid: Executing shell command with user input
func shellExec(command string) error {
	cmd := exec.Command("sh", "-c", command) // MATCH /potential command injection: variable argument passed to exec.Command/
	return cmd.Run()
}

// Invalid: CommandContext with variable argument
func runWithContext(ctx context.Context, cmd string) error {
	return exec.CommandContext(ctx, cmd).Run() // MATCH /potential command injection: variable argument passed to exec.CommandContext/
}

// Invalid: syscall.Exec with variable
func execSyscall(path string) error {
	return syscall.Exec(path, nil, nil) // MATCH /potential command injection: variable argument passed to syscall.Exec/
}

// Invalid: syscall.ForkExec with variable
func forkExecSyscall(path string) (int, error) {
	return syscall.ForkExec(path, nil, nil) // MATCH /potential command injection: variable argument passed to syscall.ForkExec/
}

// Invalid: syscall.StartProcess with variable
func startProcess(path string) (int, uintptr, error) {
	return syscall.StartProcess(path, nil, nil) // MATCH /potential command injection: variable argument passed to syscall.StartProcess/
}

// Valid: All arguments are constants
func listFiles() ([]byte, error) {
	return exec.Command("ls", "-la", "/tmp").Output()
}

// Valid: Constant command and arguments
func runKnownCommand() error {
	cmd := exec.Command("/usr/bin/git", "status")
	return cmd.Run()
}

// Valid: CommandContext with all constants
func runKnownWithContext(ctx context.Context) error {
	return exec.CommandContext(ctx, "/usr/bin/git", "status").Run()
}

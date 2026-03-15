---
title: noVariableCommandExecution
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noVariableCommandExecution`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noVariableCommandExecution:
    # rule options here
```

## Details

Audits the use of command execution functions with variable arguments.

This rule detects calls to `os/exec.Command`, `os/exec.CommandContext`, `syscall.Exec`, `syscall.ForkExec`, and `syscall.StartProcess` where the command name or arguments are derived from variables rather than constants. When user-controlled data is passed to command execution functions, it can lead to command injection attacks where an attacker executes arbitrary system commands.

The rule checks whether arguments to these functions can be resolved to constant values. If any argument is a variable that cannot be resolved to a compile-time constant, the rule flags it as a potential security issue. Parameters and struct field accesses are given special consideration in certain positions.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "os/exec"

func runCommand(userInput string) error {
    // Command name from variable
    cmd := exec.Command(userInput)
    return cmd.Run()
}
```

```golang
import "os/exec"

func processFile(filename string) ([]byte, error) {
    // Argument from variable - potential command injection
    return exec.Command("cat", filename).Output()
}
```

```golang
import "os/exec"

func shellExec(command string) error {
    // Executing shell command with user input
    cmd := exec.Command("sh", "-c", command)
    return cmd.Run()
}
```

### Valid

```golang
import "os/exec"

func listFiles() ([]byte, error) {
    // All arguments are constants
    return exec.Command("ls", "-la", "/tmp").Output()
}
```

```golang
import "os/exec"

func runKnownCommand() error {
    // Constant command and arguments
    cmd := exec.Command("/usr/bin/git", "status")
    return cmd.Run()
}
```

```golang
import (
    "os/exec"
    "path/filepath"
    "regexp"
)

func processFile(filename string) ([]byte, error) {
    // Validate input before using in command
    if !regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`).MatchString(filename) {
        return nil, fmt.Errorf("invalid filename")
    }
    cleanPath := filepath.Clean(filepath.Join("/safe/dir", filename))
    return exec.Command("cat", cleanPath).Output()
}
```

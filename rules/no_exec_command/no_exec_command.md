---
title: noExecCommand
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noExecCommand`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noExecCommand:
    # no additional options
```

## Details

Disallow calling `os/exec.Command` without a context. The function `exec.Command` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation and deadlines of the spawned process. Use `os/exec.CommandContext` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to command execution enables the caller to kill the process when the context is canceled, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import (
    "fmt"
    "os/exec"
)

func main() {
    // exec.Command does not accept a context
    cmd := exec.Command("ls", "-la")
    output, err := cmd.Output()
    if err != nil {
        panic(err)
    }
    fmt.Println(string(output))
}
```

### Valid

```golang
package main

import (
    "context"
    "fmt"
    "os/exec"
)

func main() {
    ctx := context.Background()
    cmd := exec.CommandContext(ctx, "ls", "-la")
    output, err := cmd.Output()
    if err != nil {
        panic(err)
    }
    fmt.Println(string(output))
}
```

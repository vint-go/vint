---
title: noInvalidExecCommandArg
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInvalidExecCommandArg`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInvalidExecCommandArg:
    # rule options here
```

## Details

Invalid first argument to `exec.Command`.

The first argument to `exec.Command` is the path to the executable, not a shell command. Passing a string containing spaces (like `"echo hello"`) means it will look for an executable literally named `"echo hello"` rather than running `echo` with argument `hello`.

Source: https://staticcheck.dev/docs/checks/#SA1005

## Examples

### Invalid

```golang
package main

import "os/exec"

func main() {
    // Wrong: treats "echo hello" as the executable name
    cmd := exec.Command("echo hello")
    _ = cmd
}
```

### Valid

```golang
package main

import "os/exec"

func main() {
    // Correct: "echo" is the executable, "hello" is the argument
    cmd := exec.Command("echo", "hello")
    _ = cmd
}
```

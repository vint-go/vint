---
title: noUntrappableSignal
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUntrappableSignal`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUntrappableSignal:
    # rule options here
```

## Details

Trapping a signal that cannot be intercepted.

Certain signals, such as `SIGKILL` and `SIGSTOP`, cannot be caught or ignored by a program. Attempting to trap them with `signal.Notify` or `signal.Ignore` has no effect and indicates a misunderstanding of signal handling.

Source: https://staticcheck.dev/docs/checks/#SA1016

## Examples

### Invalid

```golang
package main

import (
    "os"
    "os/signal"
    "syscall"
)

func main() {
    c := make(chan os.Signal, 1)
    // SIGKILL cannot be trapped
    signal.Notify(c, syscall.SIGKILL)
}
```

### Valid

```golang
package main

import (
    "os"
    "os/signal"
    "syscall"
)

func main() {
    c := make(chan os.Signal, 1)
    // SIGTERM can be trapped
    signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
}
```

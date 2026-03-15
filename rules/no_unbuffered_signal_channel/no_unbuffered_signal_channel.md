---
title: noUnbufferedSignalChannel
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnbufferedSignalChannel`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnbufferedSignalChannel:
    # rule options here
```

## Details

Detects misuse of unbuffered `os.Signal` channels as arguments to `signal.Notify`. The `signal.Notify` function sends signals on the provided channel in a non-blocking manner. If the channel is unbuffered (or its buffer is full), the signal will be dropped. This is almost always a bug because the programmer expects to receive all signals.

The channel passed to `signal.Notify` should be buffered with at least a capacity of 1.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/sigchanyzer

## Examples

### Invalid

```golang
import (
    "os"
    "os/signal"
    "syscall"
)

func example() {
    // Bad: unbuffered channel may miss signals
    c := make(chan os.Signal)
    signal.Notify(c, syscall.SIGINT)
}
```

### Valid

```golang
import (
    "os"
    "os/signal"
    "syscall"
)

func example() {
    // Good: buffered channel with capacity of 1
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT)
}
```

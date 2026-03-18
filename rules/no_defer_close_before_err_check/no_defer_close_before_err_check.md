---
title: noDeferCloseBeforeErrCheck
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noDeferCloseBeforeErrCheck`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noDeferCloseBeforeErrCheck:
    # rule options here
```

## Details

Deferring `Close` before checking for a possible error.

When opening a resource that implements `io.Closer`, you should check the error from the `Open` call before deferring `Close`. If `Open` returns an error, the returned value may be nil, and calling `Close` on nil will panic.

Source: https://staticcheck.dev/docs/checks/#SA5001

## Examples

### Invalid

```golang
package main

import "os"

func process() {
    f, err := os.Open("file.txt")
    // Close deferred before checking error
    defer f.Close()
    if err != nil {
        return
    }
}
```

### Valid

```golang
package main

import "os"

func process() {
    f, err := os.Open("file.txt")
    if err != nil {
        return
    }
    defer f.Close()
}
```

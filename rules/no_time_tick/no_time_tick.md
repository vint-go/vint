---
title: noTimeTick
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noTimeTick`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noTimeTick:
    # rule options here
```

## Details

Using `time.Tick` in a function other than `main()` or `init()` leads to a leak.

`time.Tick` creates a `time.Ticker` that can never be stopped or garbage collected. This is only appropriate in long-lived functions like `main()` or `init()`. In other functions, use `time.NewTicker` and call `Stop()` when done.

Source: https://staticcheck.dev/docs/checks/#SA1015

## Examples

### Invalid

```golang
package main

import "time"

func process() {
    // Tick leaks the underlying ticker
    for range time.Tick(1 * time.Second) {
        // do work
    }
}
```

### Valid

```golang
package main

import "time"

func process() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
        // do work
    }
}
```

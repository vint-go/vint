---
title: noSelectBreakConfusion
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSelectBreakConfusion`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSelectBreakConfusion:
    # rule options here
```

## Details

`for { select { ...` with unconditional break exits the select, not the loop.

A `break` inside a `select` that is a direct child of a `for` loop body only breaks out of the `select`, not the loop. Use a labeled break to exit the loop.

Only `select` statements that are direct children of the `for` body are checked. Deeply nested `select` statements (e.g., inside an `if` or another construct within the `for` body) are not flagged, as the intent is typically clear in those cases.

Source: https://staticcheck.dev/docs/checks/#SA5004

## Examples

### Invalid

```golang
package main

func process(done chan struct{}) {
    for {
        select {
        case <-done:
            // Only breaks out of select, not the loop
            break
        }
    }
}
```

### Valid

```golang
package main

func process(done chan struct{}) {
loop:
    for {
        select {
        case <-done:
            break loop
        }
    }
}
```

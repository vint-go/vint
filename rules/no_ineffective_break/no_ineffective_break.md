---
title: noIneffectiveBreak
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noIneffectiveBreak`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noIneffectiveBreak:
    # rule options here
```

## Details

Break statement with no effect. Did you mean to break out of an outer loop?

A `break` statement inside a `switch` or `select` only breaks out of the `switch`/`select`, not an enclosing `for` loop. If you intended to break out of the loop, use a labeled break.

Source: https://staticcheck.dev/docs/checks/#SA4011

## Examples

### Invalid

```golang
package main

func process(items []string) {
    for _, item := range items {
        switch item {
        case "stop":
            // Only breaks out of switch, not the for loop
            break
        }
    }
}
```

### Valid

```golang
package main

func process(items []string) {
loop:
    for _, item := range items {
        switch item {
        case "stop":
            // Breaks out of the for loop
            break loop
        }
    }
}
```

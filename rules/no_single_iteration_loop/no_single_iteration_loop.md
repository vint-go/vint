---
title: noSingleIterationLoop
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noSingleIterationLoop`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noSingleIterationLoop:
    # rule options here
```

## Details

The loop exits unconditionally after one iteration.

A loop that contains only a `break` or `return` in its body will never iterate more than once, making the loop pointless.

Source: https://staticcheck.dev/docs/checks/#SA4004

## Examples

### Invalid

```golang
package main

func firstItem(items []string) string {
    for _, item := range items {
        // Loop always exits on first iteration
        return item
    }
    return ""
}
```

### Valid

```golang
package main

func firstItem(items []string) string {
    if len(items) > 0 {
        return items[0]
    }
    return ""
}
```

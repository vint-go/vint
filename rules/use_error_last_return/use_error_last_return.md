---
title: useErrorLastReturn
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useErrorLastReturn`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useErrorLastReturn:
    # rule options here
```

## Details

A function's error value should be its last return value.

By Go convention, functions that return an error should return it as the last value. This is a widely followed convention that makes error handling patterns consistent.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
package main

// Error is not the last return value
func process() (error, string) {
    return nil, "ok"
}
```

### Valid

```golang
package main

// Error is the last return value
func process() (string, error) {
    return "ok", nil
}
```

---
title: useInfiniteFor
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useInfiniteFor`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useInfiniteFor:
    # no additional options
```

## Details

Use `for { ... }` instead of `for true { ... }`.

The `for true` construct is equivalent to `for` (an unconditional loop). The idiomatic way in Go is to omit the condition entirely.

Source: https://staticcheck.dev/docs/checks/#S1006

## Examples

### Invalid

```golang
package main

func process() {
    for true {
        // ...
        break
    }
}
```

### Valid

```golang
package main

func process() {
    for {
        // ...
        break
    }
}
```

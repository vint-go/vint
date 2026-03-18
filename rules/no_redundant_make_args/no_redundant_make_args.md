---
title: noRedundantMakeArgs
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantMakeArgs`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantMakeArgs:
    # rule options here
```

## Details

Simplify `make` call by omitting redundant arguments.

For slices, `make([]T, N, N)` where length and capacity are the same can be simplified to `make([]T, N)`. For maps, `make(map[K]V, 0)` can be simplified to `make(map[K]V)`.

Source: https://staticcheck.dev/docs/checks/#S1019

## Examples

### Invalid

```golang
package main

func process() {
    // Redundant capacity argument
    s := make([]int, 10, 10)
    _ = s
}
```

### Valid

```golang
package main

func process() {
    s := make([]int, 10)
    _ = s
}
```

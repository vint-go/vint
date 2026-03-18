---
title: useMergedVarDecl
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useMergedVarDecl`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useMergedVarDecl:
    # rule options here
```

## Details

Merge variable declaration and assignment.

A variable declaration followed immediately by an assignment can be combined into a short variable declaration.

Source: https://staticcheck.dev/docs/checks/#S1021

## Examples

### Invalid

```golang
package main

func process() {
    var x int
    x = 42
    _ = x
}
```

### Valid

```golang
package main

func process() {
    x := 42
    _ = x
}
```

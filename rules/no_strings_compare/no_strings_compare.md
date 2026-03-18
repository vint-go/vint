---
title: noStringsCompare
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noStringsCompare`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noStringsCompare:
    # no additional options
```

## Details

Detects usage of `strings.Compare`. The Go documentation itself recommends against using `strings.Compare`, stating that it exists only for symmetry with package `bytes`. Using `==`, `<`, and `>` operators directly is clearer and more efficient for string comparison.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
if strings.Compare(a, b) == 0 {
    // equal
}
```

```golang
if strings.Compare(a, b) < 0 {
    // a < b
}
```

### Valid

```golang
if a == b {
    // equal
}
```

```golang
if a < b {
    // a < b
}
```

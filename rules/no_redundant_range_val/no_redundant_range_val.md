---
title: noRedundantRangeVal
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantRangeVal`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantRangeVal:
    # no options
```

## Details

This rule suggests a shorter way of writing ranges that do not use the second value. When iterating over a collection and assigning the second value to the blank identifier `_`, the second variable should be omitted entirely for clarity.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
for k, _ := range m {
    _ = k
}
```

```golang
for i, _ := range s {
    _ = i
}
```

### Valid

```golang
// Using only the key (single variable form)
for k := range m {
    _ = k
}
```

```golang
// Using both key and value
for k, v := range m {
    _ = k
    _ = v
}
```

```golang
// Using only the value with blank key
for _, v := range s {
    _ = v
}
```

---
title: noSuspiciousMapKey
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noSuspiciousMapKey`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noSuspiciousMapKey:
    # no additional options
```

## Details

Detects suspicious map literal keys. This checker identifies two types of issues in string-keyed map literals:

1. **Whitespace anomalies** -- Flags keys with leading or trailing single spaces, which may indicate accidental padding (e.g., `"bar "` instead of `"bar"`).
2. **Duplicate keys** -- Warns about repeated keys in the same map literal for non-basic literals that don't have side effects.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
m := map[string]int{
    "foo":  1,
    "bar ": 2, // trailing space in key
}
```

```golang
m := map[string]int{
    getValue(): 1,
    getValue(): 2, // duplicate non-literal key
}
```

### Valid

```golang
m := map[string]int{
    "foo": 1,
    "bar": 2,
}
```

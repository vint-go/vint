---
title: noRangeAppendAll
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noRangeAppendAll`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noRangeAppendAll:
    # no additional options
```

## Details

Detects appending an entire slice inside a range loop over that same slice. This is a performance anti-pattern where code appends the full slice using the ellipsis operator (`...`) while simultaneously iterating over it. Each iteration copies the entire slice again, leading to quadratic behavior.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
for _, n := range ns {
    rs = append(rs, ns...) // appending all of ns on each iteration
}
```

### Valid

```golang
for _, n := range ns {
    rs = append(rs, n) // append just the current element
}
```

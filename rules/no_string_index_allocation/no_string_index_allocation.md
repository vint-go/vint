---
title: noStringIndexAllocation
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/noStringIndexAllocation`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/noStringIndexAllocation:
    # no additional options
```

## Details

Detects `strings.Index` calls that may cause unwanted allocations. When calling `strings.Index` with a `[]byte` converted to string, it causes an allocation. Using `bytes.Index` directly avoids this allocation.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
i := strings.Index(string(b), "pattern") // allocates a new string
```

### Valid

```golang
i := bytes.Index(b, []byte("pattern")) // no string allocation
```

---
title: noOffByOneError
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noOffByOneError`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noOffByOneError:
    # no additional options
```

## Details

Detects various off-by-one kind errors. This checker identifies patterns where an index or boundary calculation is off by one, which is a common source of bugs in programming. It checks for cases like accessing `s[len(s)]` instead of `s[len(s)-1]`.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// Accessing one past the last element
last := s[len(s)]
```

### Valid

```golang
// Correct last element access
last := s[len(s)-1]
```

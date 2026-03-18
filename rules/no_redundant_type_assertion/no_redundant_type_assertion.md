---
title: noRedundantTypeAssertion
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantTypeAssertion`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantTypeAssertion:
    # no additional options
```

## Details

Detects redundant type assertions. This checker identifies type assertions where the source and destination types are identical. When a value is cast to a type it already possesses, the assertion is unnecessary and can be removed.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
var r io.Reader = getReader()
v := r.(io.Reader) // redundant: r is already io.Reader
```

### Valid

```golang
var r io.Reader = getReader()
v := r // no assertion needed, already the correct type
```

```golang
var r io.Reader = getReader()
v := r.(io.ReadCloser) // valid: asserting to a different type
```

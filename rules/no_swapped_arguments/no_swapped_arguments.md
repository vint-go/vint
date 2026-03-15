---
title: noSwappedArguments
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noSwappedArguments`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noSwappedArguments:
    # no additional options
```

## Details

Detects suspicious arguments order. This checker identifies function calls where the argument order appears to be swapped, which is a common mistake. For example, passing the arguments to `strings.HasPrefix` or `strings.HasSuffix` in the wrong order.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// Arguments might be in the wrong order
strings.HasPrefix("prefix", s)
```

### Valid

```golang
strings.HasPrefix(s, "prefix")
```

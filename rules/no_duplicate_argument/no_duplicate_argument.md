---
title: noDuplicateArgument
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noDuplicateArgument`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noDuplicateArgument:
    # no additional options
```

## Details

Detects suspicious duplicated arguments. This checker identifies function calls where the same argument is passed twice when different arguments were likely intended. This is a common copy-paste error.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
copy(dst, dst) // likely meant copy(dst, src)
```

```golang
reflect.DeepEqual(x, x) // comparing x to itself is always true
```

### Valid

```golang
copy(dst, src)
```

```golang
reflect.DeepEqual(x, y)
```

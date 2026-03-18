---
title: noAlwaysTrueLenCheck
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noAlwaysTrueLenCheck`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noAlwaysTrueLenCheck:
    # no additional options
```

## Details

Detects usage of `len` when the result is obvious or can be replaced with a more idiomatic check. For example, checking `len(s) >= 0` is always true for any slice or string, and `len(s) < 0` is always false. These are likely bugs or unnecessary checks.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
if len(s) >= 0 { // always true
    doSomething()
}
```

```golang
if len(s) < 0 { // always false
    doSomething()
}
```

### Valid

```golang
if len(s) > 0 {
    doSomething()
}
```

```golang
if len(s) == 0 {
    handleEmpty()
}
```

---
title: noUncheckedInlineError
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noUncheckedInlineError`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noUncheckedInlineError:
    # no additional options
```

## Details

Detects unchecked errors in if statement assignments. This checker identifies patterns where an error value assigned in an if statement's initialization clause is not properly checked in the condition. The error might be silently ignored, leading to unexpected behavior.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
if val, err := doSomething(); val != nil {
    // err is not checked
    use(val)
}
```

### Valid

```golang
if val, err := doSomething(); err == nil {
    use(val)
}
```

```golang
val, err := doSomething()
if err != nil {
    return err
}
use(val)
```

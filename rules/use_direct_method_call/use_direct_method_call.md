---
title: useDirectMethodCall
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useDirectMethodCall`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useDirectMethodCall:
    # no additional options
```

## Details

Detects method expression call that can be replaced with a method call. This checker identifies instances where a method is invoked using method expression syntax (passing the receiver as the first argument) when a direct method call would be more idiomatic.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// Method expression call with explicit receiver
foo.bar(f)           // can be replaced with f.bar()
```

### Valid

```golang
// Direct method call
f.bar()
```

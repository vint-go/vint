---
title: noDynamicFormatString
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noDynamicFormatString`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noDynamicFormatString:
    # no additional options
```

## Details

Detects suspicious formatting string usage where a dynamic (non-literal) string is passed as the format argument to `fmt.Sprintf`, `fmt.Printf`, and similar functions without any additional arguments. This is often a mistake where `fmt.Sprint` or direct string usage should be used instead, and can be a security concern if the dynamic string contains `%` characters that would be interpreted as format verbs.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
msg := fmt.Sprintf(dynamicString) // no format args, dynamic string may contain %
```

### Valid

```golang
msg := dynamicString // use directly
```

```golang
msg := fmt.Sprintf("%s", dynamicString) // explicit format verb
```

```golang
msg := fmt.Sprintf("prefix: %s", dynamicString) // proper formatting
```

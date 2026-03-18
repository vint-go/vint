---
title: usePercentQ
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/usePercentQ`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/usePercentQ:
    # no additional options
```

## Details

Detects `"%s"` formatting that can be replaced with `%q`. When a string is formatted with `fmt.Sprintf` using `"%s"` within surrounding quotes, the `%q` verb should be used instead. The `%q` verb automatically adds double quotes and properly escapes special characters.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
msg := fmt.Sprintf("value is \"%s\"", s)
```

### Valid

```golang
msg := fmt.Sprintf("value is %q", s)
```

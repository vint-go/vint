---
title: useEqualFold
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/useEqualFold`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/useEqualFold:
    # no additional options
```

## Details

Detects case-insensitive string comparisons that convert strings to lower or upper case before comparing, instead of using `strings.EqualFold`. Converting with `strings.ToLower` or `strings.ToUpper` allocates a new string, while `strings.EqualFold` performs the comparison in place without allocation and handles Unicode case folding correctly.

The rule only flags comparisons where both operands are **pure** (side-effect free). Expressions involving function or method calls are considered impure and are skipped, since replacing them with `strings.EqualFold` could change evaluation semantics. Self-comparisons (both sides textually identical) are also skipped.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
if strings.ToLower(a) == strings.ToLower(b) {
    // case-insensitive comparison via ToLower
}
```

```golang
if strings.ToUpper(s) == "HELLO" {
    // suboptimal case-insensitive comparison
}
```

### Valid

```golang
if strings.EqualFold(a, b) {
    // efficient case-insensitive comparison
}
```

```golang
if strings.EqualFold(s, "HELLO") {
    // correct and efficient
}
```

```golang
// Not flagged: function call argument is impure
if strings.ToLower(getConfig("env")) == "production" {
    // ...
}
```

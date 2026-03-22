---
title: noNonWrappingErrorf
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNonWrappingErrorf`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNonWrappingErrorf:
    # whether to allow multiple %w verbs (Go 1.20+)
    errorf-multi: true
```

## Details

Detects `fmt.Errorf` calls that use non-wrapping format verbs (`%v`, `%s`, etc.) instead of the `%w` verb for error values. Go 1.13 introduced error wrapping with `%w`, which preserves the error chain and allows callers to use `errors.Is` and `errors.As` to inspect wrapped errors. Using `%v` instead of `%w` loses this error chain information.

Note that wrapping an error makes it part of your API surface. If you intentionally want to prevent callers from matching the error, using `%v` is appropriate.

Starting with Go 1.20, multiple `%w` verbs are supported in a single `fmt.Errorf` call.

Source: https://github.com/polyfloyd/go-errorlint

## Examples

### Invalid

```golang
// Using %v loses the error chain
fmt.Errorf("failed to process: %v", err)
```

```golang
// Using %s also loses the error chain
fmt.Errorf("operation failed: %s", err)
```

### Valid

```golang
// Using %w preserves the error chain
fmt.Errorf("failed to process: %w", err)
```

```golang
// Multiple %w verbs (Go 1.20+)
fmt.Errorf("error1: %w, error2: %w", err1, err2)
```

```golang
// Intentionally not wrapping to hide internal error
fmt.Errorf("operation failed: %v", internalErr)
```

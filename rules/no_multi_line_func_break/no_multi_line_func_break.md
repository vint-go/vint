---
title: noMultiLineFuncBreak
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noMultiLineFuncBreak`
- This rule is not recommended, meaning it is **not** enabled by default.
- This rule has a **fix**.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noMultiLineFuncBreak:
    multi-func: true # Enable the multi-line function signature check. Default: false.
```

## Details

Detects multi-line function signatures where the parameter list does not start on the same line as the `func` keyword. When a function has many parameters and the signature spans multiple lines, the first parameter should appear on the same line as `func`, rather than starting on a new line below it.

This enforces a consistent style for multi-line function declarations and improves readability. Without this check, a developer might place the `func` keyword and function name on their own line, with all parameters starting on the next line, which wastes vertical space and can reduce clarity about where the signature begins.

This check is **not enabled by default**. Set the `multi-func` option to `true` to activate it.

Source: https://github.com/ultraware/whitespace

## Examples

### Invalid

```golang
// The parameter list starts on a new line after the function name.
func createUser(
	name string,
	email string,
	age int,
) (*User, error) {
	return &User{Name: name, Email: email, Age: age}, nil
}
```

```golang
// The parameter list starts on a new line after the function name.
func (s *Service) ProcessOrder(
	ctx context.Context,
	orderID string,
	options ...Option,
) error {
	return s.process(ctx, orderID, options...)
}
```

### Valid

```golang
// The first parameter is on the same line as the function name.
func createUser(name string,
	email string,
	age int,
) (*User, error) {
	return &User{Name: name, Email: email, Age: age}, nil
}
```

```golang
// The first parameter is on the same line as the function name.
func (s *Service) ProcessOrder(ctx context.Context,
	orderID string,
	options ...Option,
) error {
	return s.process(ctx, orderID, options...)
}
```

```golang
// Single-line function signature; no issue.
func add(a, b int) int {
	return a + b
}
```

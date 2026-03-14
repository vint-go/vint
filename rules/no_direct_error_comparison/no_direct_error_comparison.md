---
title: noDirectErrorComparison
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noDirectErrorComparison`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has an auto-fix that replaces the comparison with `errors.Is()`.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noDirectErrorComparison:
    # no additional options
```

## Details

This rule flags direct comparisons of error values using `==` or `!=` operators and recommends using `errors.Is()` instead. Since Go 1.13, errors can be wrapped using `fmt.Errorf` with the `%w` verb. When errors are wrapped, a direct comparison with `==` will fail to match the underlying error, because the wrapped error is a different value. The `errors.Is()` function traverses the chain of wrapped errors, ensuring correct matching regardless of wrapping.

The following comparisons are allowed and will not trigger the rule:
- Comparing an error to `nil` (e.g., `if err != nil`)
- Comparing an error to `io.EOF` (a well-known sentinel that is conventionally compared directly)

When this rule triggers, it provides a suggested fix that rewrites the comparison to use `errors.Is()`. For example, `err == ErrFoo` would be rewritten to `errors.Is(err, ErrFoo)`, and `err != ErrFoo` would be rewritten to `!errors.Is(err, ErrFoo)`.

Source: https://github.com/Djarvur/go-err113

## Examples

### Invalid

```golang
// Direct equality comparison of errors
func handleError(err error) {
	if err == ErrPermission {
		// handle permission error
	}
}
```

```golang
// Direct inequality comparison of errors
func handleError(err error) {
	if err != ErrTimeout {
		// handle non-timeout error
	}
}
```

```golang
// Comparing errors returned from function calls
func process() {
	err := doSomething()
	if err == getExpectedError() {
		return
	}
}
```

### Valid

```golang
// Using errors.Is() for comparison
func handleError(err error) {
	if errors.Is(err, ErrPermission) {
		// handle permission error
	}
}
```

```golang
// Comparing to nil is allowed
func handleError(err error) {
	if err != nil {
		// handle error
	}
}
```

```golang
// Comparing to io.EOF is allowed
func readAll(r io.Reader) {
	_, err := r.Read(buf)
	if err == io.EOF {
		return
	}
}
```

```golang
// Using !errors.Is() for inequality
func handleError(err error) {
	if !errors.Is(err, ErrTimeout) {
		// handle non-timeout error
	}
}
```

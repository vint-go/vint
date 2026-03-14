---
title: noDynamicErrors
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noDynamicErrors`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noDynamicErrors:
    # no additional options
```

## Details

This rule flags the creation of dynamic errors inside functions and requires that errors be defined as package-level variables (sentinel errors) and wrapped when additional context is needed. Since Go 1.13, the idiomatic approach is to define static errors at the package level using `errors.New()` or `var ErrFoo = errors.New("foo")`, and then wrap them with `fmt.Errorf("context: %w", ErrFoo)` when returning from functions.

Creating errors dynamically within function bodies using `errors.New()` or `fmt.Errorf()` (without the `%w` wrapping verb) makes it impossible for callers to reliably check for specific error conditions using `errors.Is()` or `errors.As()`, because each call creates a distinct error value.

The rule specifically checks for:
- **`errors.New()` calls inside function bodies** -- these are always flagged, since they produce a new unique error value each time.
- **`fmt.Errorf()` calls inside function bodies without `%w`** -- these are flagged because they create a new error without wrapping an existing sentinel. If the format string contains `%w`, the call is allowed because it wraps an existing error.

The following patterns are allowed and will not trigger the rule:
- Package-level variable declarations using `errors.New()` (e.g., `var ErrFoo = errors.New("foo")`)
- `fmt.Errorf()` calls that use the `%w` verb to wrap an existing error
- Error creation from non-standard error packages (e.g., `github.com/pkg/errors`) is not checked

Source: https://github.com/Djarvur/go-err113

## Examples

### Invalid

```golang
// Creating a dynamic error with errors.New() inside a function
func validate(name string) error {
	if name == "" {
		return errors.New("name is required")
	}
	return nil
}
```

```golang
// Creating a dynamic error with fmt.Errorf() without wrapping
func openFile(path string) error {
	_, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", path, err)
	}
	return nil
}
```

```golang
// Using errors.New() in a return statement
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}
```

### Valid

```golang
// Defining sentinel errors at the package level
var (
	ErrNameRequired  = errors.New("name is required")
	ErrDivisionByZero = errors.New("division by zero")
)

func validate(name string) error {
	if name == "" {
		return ErrNameRequired
	}
	return nil
}
```

```golang
// Wrapping a sentinel error with fmt.Errorf using %w
var ErrOpenFile = errors.New("failed to open file")

func openFile(path string) error {
	_, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	return nil
}
```

```golang
// Package-level error definitions are allowed
var ErrNotFound = errors.New("not found")
var ErrUnauthorized = errors.New("unauthorized")
```

```golang
// Wrapping errors with %w is allowed inside functions
func process(id int) error {
	item, err := fetch(id)
	if err != nil {
		return fmt.Errorf("processing item %d: %w", id, err)
	}
	return nil
}
```

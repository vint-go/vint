---
title: useErrorsAs
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/useErrorsAs`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/useErrorsAs:
    # no additional options
```

## Details

Detects type assertions and type switches on error values that should use `errors.As()` instead. Type assertions like `err.(*MyError)` fail when errors are wrapped, because the outer wrapper type does not match the inner error type.

Go 1.13 introduced `errors.As()` which traverses the error chain and finds the first error that matches the target type, making it the correct way to extract a specific error type from a chain.

Source: https://github.com/polyfloyd/go-errorlint

## Examples

### Invalid

```go
// Type assertion misses wrapped errors
myErr, ok := err.(*MyError)
if ok {
    fmt.Println(myErr.Code)
}
```

```go
// Type switch misses wrapped errors
switch e := err.(type) {
case *MyError:
    fmt.Println(e.Code)
case *OtherError:
    fmt.Println(e.Message)
}
```

### Valid

```go
// Using errors.As traverses the error chain
var myErr *MyError
if errors.As(err, &myErr) {
    fmt.Println(myErr.Code)
}
```

```go
// Check for multiple error types with errors.As
var myErr *MyError
var otherErr *OtherError
if errors.As(err, &myErr) {
    fmt.Println(myErr.Code)
} else if errors.As(err, &otherErr) {
    fmt.Println(otherErr.Message)
}
```

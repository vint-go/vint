---
title: noInvalidErrorsAs
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInvalidErrorsAs`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInvalidErrorsAs:
    # rule options here
```

## Details

Checks that the second argument to `errors.As` is a pointer to a type implementing the `error` interface, or a pointer to any interface type. The `errors.As` function requires its second argument to be a non-nil pointer to either a type that implements error, or to any interface type. Passing a value that does not meet these requirements will cause a runtime panic.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/errorsas

## Examples

### Invalid

```golang
import "errors"

type MyError struct {
    Code int
}

func (e *MyError) Error() string {
    return "my error"
}

func example(err error) {
    var myErr MyError
    // Bad: should pass a pointer to *MyError, not a pointer to MyError
    if errors.As(err, &myErr) {
        // ...
    }
}
```

### Valid

```golang
import "errors"

type MyError struct {
    Code int
}

func (e *MyError) Error() string {
    return "my error"
}

func example(err error) {
    var myErr *MyError
    // Good: passing a pointer to *MyError
    if errors.As(err, &myErr) {
        fmt.Println(myErr.Code)
    }
}
```

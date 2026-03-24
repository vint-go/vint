---
title: noImpossibleNilComparison
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noImpossibleNilComparison`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noImpossibleNilComparison:
    # rule options here
```

## Details

Detects nil comparisons of interface values that are always false (`== nil`) or always true (`!= nil`).

This rule checks at the **call site**, matching the semantics of staticcheck SA4023. It flags `x == nil` or `x != nil` comparisons where static analysis proves the interface value can never be untyped nil. This happens in two cases:

1. The value was assigned from a function that provably never returns untyped nil at that return position (all return paths return concrete types).
2. The value was assigned from a function whose declared return type is a concrete (non-interface) type, meaning the result is always wrapped in an interface and can never be nil.

An interface value is only nil when both its type and value components are nil (untyped nil). When a concrete type (even a nil pointer) is stored in an interface, the interface itself is non-nil.

Source: https://staticcheck.dev/docs/checks/#SA4023

## Examples

### Invalid

```golang
package main

import "fmt"

type MyError struct{}

func (e *MyError) Error() string { return "error" }

// getError never returns untyped nil - all paths return a concrete type.
func getError() error {
    var err *MyError
    return err // returns non-nil interface wrapping a nil *MyError
}

func main() {
    err := getError()
    // Flagged: nil comparison of err is always false because
    // getError never returns nil.
    if err == nil {
        fmt.Println("no error") // never reached
    }
}
```

### Valid

```golang
package main

import "fmt"

type MyError struct{}

func (e *MyError) Error() string { return "error" }

// getError can return untyped nil, so nil comparisons are valid.
func getError() error {
    var err *MyError
    if err == nil {
        return nil // returns untyped nil - interface will be nil
    }
    return err
}

func main() {
    err := getError()
    if err == nil {
        fmt.Println("no error") // can be reached
    }
}
```

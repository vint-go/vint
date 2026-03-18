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

Impossible comparison of interface value with untyped nil.

When a concrete type is stored in an interface, comparing the interface to `nil` will be false even if the concrete value is the zero value. An interface is only nil when both its type and value are nil.

Source: https://staticcheck.dev/docs/checks/#SA4023

## Examples

### Invalid

```golang
package main

import "fmt"

type MyError struct{}

func (e *MyError) Error() string { return "error" }

func getError() error {
    var err *MyError
    // Returns non-nil interface even though err is nil
    return err
}

func main() {
    err := getError()
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

func getError() error {
    var err *MyError
    if err == nil {
        return nil // Return untyped nil
    }
    return err
}

func main() {
    err := getError()
    if err == nil {
        fmt.Println("no error")
    }
}
```

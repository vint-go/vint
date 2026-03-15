---
title: noDeepEqualErrors
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noDeepEqualErrors`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noDeepEqualErrors:
    # rule options here
```

## Details

Checks for the use of `reflect.DeepEqual` with error values. Using `reflect.DeepEqual` to compare errors is problematic because error types often have unexported fields or are interface values. Two errors that represent the same condition might not be deeply equal, and two errors that appear equal might not be semantically equivalent.

Instead, use `errors.Is` or compare specific error properties to check error equality.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/deepequalerrors

## Examples

### Invalid

```golang
import (
    "errors"
    "reflect"
)

func example() {
    err1 := errors.New("something failed")
    err2 := errors.New("something failed")
    // Bad: using reflect.DeepEqual with errors
    if reflect.DeepEqual(err1, err2) {
        // ...
    }
}
```

### Valid

```golang
import (
    "errors"
    "io"
)

func example(err error) {
    // Good: use errors.Is for error comparison
    if errors.Is(err, io.EOF) {
        // ...
    }
}
```

```golang
import "errors"

var ErrNotFound = errors.New("not found")

func example(err error) {
    // Good: use errors.Is for sentinel error comparison
    if errors.Is(err, ErrNotFound) {
        // ...
    }
}
```

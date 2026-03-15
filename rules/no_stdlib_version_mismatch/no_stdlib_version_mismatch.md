---
title: noStdlibVersionMismatch
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noStdlibVersionMismatch`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noStdlibVersionMismatch:
    # rule options here
```

## Details

Reports uses of standard library symbols that are "too new" for the Go version in effect. If your `go.mod` file specifies `go 1.20` but your code uses a function introduced in Go 1.21, this analyzer will flag it. This helps ensure that your code is compatible with the minimum Go version declared in your module.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/stdversion

## Examples

### Invalid

```golang
// go.mod specifies: go 1.20

package example

import "slices" // Bad: slices package was added in Go 1.21

func example() {
    s := []int{3, 1, 2}
    slices.Sort(s) // Bad: slices.Sort requires Go 1.21
}
```

### Valid

```golang
// go.mod specifies: go 1.21

package example

import "slices" // Good: module requires Go 1.21 or later

func example() {
    s := []int{3, 1, 2}
    slices.Sort(s)
}
```

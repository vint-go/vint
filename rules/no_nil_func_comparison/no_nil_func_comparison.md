---
title: noNilFuncComparison
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNilFuncComparison`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNilFuncComparison:
    # rule options here
```

## Details

Checks for useless comparisons against nil. This analyzer detects comparisons of functions against nil that are always true or always false. A function value is only nil if it was never assigned, but if the comparison is against a named function (not a function variable), the result is always the same -- named functions are never nil.

This usually indicates a bug where the programmer intended to call the function rather than compare it to nil.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/nilfunc

## Examples

### Invalid

```golang
func myFunc() {}

func example() {
    // Bad: comparing a named function to nil is always false
    if myFunc == nil {
        fmt.Println("unreachable")
    }
}
```

### Valid

```golang
func example() {
    var fn func()
    // Good: comparing a function variable to nil is meaningful
    if fn == nil {
        fmt.Println("fn is not set")
    }
}
```

```golang
func myFunc() error {
    return nil
}

func example() {
    // Good: comparing the result of a function call to nil
    if myFunc() == nil {
        fmt.Println("no error")
    }
}
```

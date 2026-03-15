---
title: noSpecificFunctionCall
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noSpecificFunctionCall`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noSpecificFunctionCall:
    name: "example" # function name to search for
```

## Details

A trivial example analyzer that reports calls to a specific function. This analyzer is primarily used as a demonstration and test of the Analysis API. It finds all calls to a function with a configurable name and reports them.

While mainly educational, it can be configured to find calls to any named function, which may occasionally be useful for auditing purposes.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/findcall

## Examples

### Invalid

```golang
// If configured with name: "println"
func example() {
    println("hello") // reported: call of println(...)
}
```

### Valid

```golang
// If configured with name: "println"
func example() {
    fmt.Println("hello") // not reported: different function
}
```

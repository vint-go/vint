---
title: noVariableShadowing
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noVariableShadowing`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noVariableShadowing:
    # rule options here
```

## Details

Checks for shadowed variables. Variable shadowing occurs when a variable declared in an inner scope has the same name as a variable in an outer scope. This can lead to subtle bugs where the programmer intends to use the outer variable but accidentally creates a new one with `:=`.

This analyzer reports cases where a short variable declaration (`:=`) shadows a variable from an outer scope, which is a common source of bugs.

Note: This analyzer is not enabled by default because variable shadowing is sometimes intentional and the false positive rate can be high.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/shadow

## Examples

### Invalid

```golang
func example() (err error) {
    x, err := doFirst()
    if err != nil {
        return err
    }
    // Bad: err is shadowed by the inner := declaration
    // The outer err (named return) is not updated
    if y, err := doSecond(x); err != nil {
        return err // returns inner err, but outer err is still nil
    }
    return nil
}
```

### Valid

```golang
func example() (err error) {
    x, err := doFirst()
    if err != nil {
        return err
    }
    // Good: use = instead of := to avoid shadowing
    var y int
    y, err = doSecond(x)
    if err != nil {
        return err
    }
    _ = y
    return nil
}
```

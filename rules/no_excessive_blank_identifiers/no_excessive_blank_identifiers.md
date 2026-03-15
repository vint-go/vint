---
title: noExcessiveBlankIdentifiers
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noExcessiveBlankIdentifiers`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noExcessiveBlankIdentifiers:
    max-blank-identifiers: 2
```

## Details

Checks assignments with too many blank identifiers (e.g. `x, _, _, _, := f()`).

When a function returns multiple values, Go allows you to use the blank identifier `_` to discard unwanted return values. However, having too many blank identifiers in a single assignment statement is a code smell. It often indicates that the function's API is being misused or that the code could be restructured. A large number of discarded return values suggests the developer may not fully understand what the function returns, or that the function itself may need refactoring.

By default, this rule flags any assignment that contains more than 2 blank identifiers. This threshold is configurable via the `max-blank-identifiers` option.

Source: https://github.com/golangci/golangci-lint/blob/master/pkg/golinters/dogsled/dogsled.go

## Examples

### Invalid

```golang
// With default max-blank-identifiers: 2, this has 3 blank identifiers
a, _, _, _, b := myFunc()
```

```golang
// All return values except one are discarded
x, _, _, _ := fourReturnValues()
```

```golang
// Excessive blank identifiers in a short variable declaration
result, _, _, _, _, _ := complexFunction()
```

### Valid

```golang
// Only 1 blank identifier (within default threshold of 2)
a, _ := myFunc()
```

```golang
// Exactly 2 blank identifiers (within default threshold of 2)
a, _, _ := myFunc()
```

```golang
// No blank identifiers at all
a, b, c := myFunc()
```

```golang
// Single return value assignment
result := singleReturn()
```

```golang
// Using all return values
value, err := os.Open("file.txt")
```

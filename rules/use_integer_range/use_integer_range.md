---
title: useIntegerRange
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useIntegerRange`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useIntegerRange:
    # no additional options
```

## Details

Detects traditional C-style `for` loops that can be replaced with Go 1.22+ integer range syntax.

Go 1.22 introduced the ability to range over integers, allowing `for i := range n` as a cleaner alternative to `for i := 0; i < n; i++`. This rule identifies loops that follow the classic C-style pattern (initializing at 0, comparing against an upper bound, and incrementing by 1) and suggests converting them to the more concise integer range form.

The rule recognizes several increment styles including `i++`, `i += 1`, `i = i + 1`, and `i = 1 + i`. It also handles various comparison operators (`<`, `<=`, `>`, `>=`).

When the loop variable is not used in the body, the rule suggests using `for range n` instead.

Source: https://github.com/ckaznocha/intrange

## Examples

### Invalid

```golang
// C-style for loop that can use integer range
for i := 0; i < 10; i++ {
    fmt.Println(i)
}
```

```golang
// Unused loop variable
for i := 0; i < 10; i++ {
    fmt.Println("Hello!")
}
```

```golang
// Using += 1 increment
for i := 0; i < n; i += 1 {
    doSomething(i)
}
```

```golang
// Using i = i + 1 increment
for i := 0; i < n; i = i + 1 {
    doSomething(i)
}
```

### Valid

```golang
// Using integer range syntax (Go 1.22+)
for i := range 10 {
    fmt.Println(i)
}
```

```golang
// Using range without variable when loop variable is unused
for range 10 {
    fmt.Println("Hello!")
}
```

```golang
// Loop that modifies the loop variable in the body (cannot be converted)
for i := 0; i < 10; i++ {
    i += 2
    fmt.Println(i)
}
```

```golang
// Loop that doesn't start at 0 (cannot be converted)
for i := 1; i < 10; i++ {
    fmt.Println(i)
}
```

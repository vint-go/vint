---
title: noUnnecessaryLoopVarCopy
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noUnnecessaryLoopVarCopy`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a fix for simple single-assignment cases.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noUnnecessaryLoopVarCopy:
    check-alias: false # when true, also flags copies of loop variables to differently-named variables
```

## Details

Detects unnecessary copies of loop variables inside loop bodies. Starting with Go 1.22, the Go specification changed loop variable scoping so that each iteration of a `for` loop gets its own copy of the loop variable. Before Go 1.22, loop variables were shared across all iterations, which led to a common bug when closures or goroutines captured the variable. Developers adopted the idiom `v := v` to create a per-iteration copy.

With Go 1.22 and later, this copying idiom is no longer necessary because the language itself guarantees per-iteration scoping. Keeping these redundant copy statements adds noise to the code and may confuse readers into thinking there is still a scoping issue.

By default, this rule only flags assignments where the loop variable is copied to a new variable with the **same name** (e.g., `i := i`). When the `check-alias` option is enabled, it additionally flags assignments where the loop variable is copied to a **differently-named** variable (e.g., `_i := i`), since those are also unnecessary under Go 1.22+ semantics.

This rule applies to both `range` loops and traditional C-style `for` loops.

The diagnostic message produced is: `The copy of the 'for' variable "<name>" can be deleted (Go 1.22+)`

For simple single-assignment statements (e.g., `i := i`), an automatic fix is available that removes the entire assignment line.

Source: https://github.com/karamaru-alpha/copyloopvar

## Examples

### Invalid

```golang
// Range loop: copying both index and value variables to same names
for i, v := range []int{1, 2, 3} {
    i := i // The copy of the 'for' variable "i" can be deleted (Go 1.22+)
    v := v // The copy of the 'for' variable "v" can be deleted (Go 1.22+)
    fmt.Println(i, v)
}
```

```golang
// C-style for loop: copying loop init variable to same name
for i := 0; i < 10; i++ {
    i := i // The copy of the 'for' variable "i" can be deleted (Go 1.22+)
    go func() {
        fmt.Println(i)
    }()
}
```

```golang
// C-style for loop with multiple init variables
for i, j := 1, 1; i+j <= 3; i++ {
    i := i // The copy of the 'for' variable "i" can be deleted (Go 1.22+)
    j := j // The copy of the 'for' variable "j" can be deleted (Go 1.22+)
    fmt.Println(i, j)
}
```

```golang
// Multi-value assignment where one of the values is a loop variable copy
for i, v := range []int{1, 2, 3} {
    a, i := 1, i // The copy of the 'for' variable "i" can be deleted (Go 1.22+)
    b, v := 1, v // The copy of the 'for' variable "v" can be deleted (Go 1.22+)
    fmt.Println(a, i, b, v)
}
```

```golang
// With check-alias enabled: copying loop variable to a differently-named variable
for i, v := range []int{1, 2, 3} {
    _i := i // The copy of the 'for' variable "i" can be deleted (Go 1.22+)
    _v := v // The copy of the 'for' variable "v" can be deleted (Go 1.22+)
    fmt.Println(_i, _v)
}
```

```golang
// With check-alias enabled: multi-value assignment with aliased loop variable
for i := range []int{1, 2, 3} {
    b, _i := 1, i // The copy of the 'for' variable "i" can be deleted (Go 1.22+)
    fmt.Println(b, _i)
}
```

### Valid

```golang
// Using loop variables directly without copying
for i, v := range []int{1, 2, 3} {
    fmt.Println(i, v)
}
```

```golang
// Using loop variable directly in a C-style for loop
for i := 0; i < 10; i++ {
    go func() {
        fmt.Println(i)
    }()
}
```

```golang
// Assigning to a struct field (not a short variable declaration)
var t struct {
    Bool bool
}
for _, t.Bool = range []bool{true, false} {
    t.Bool = t.Bool // Not flagged: this is a regular assignment, not :=
}
```

```golang
// Creating a new variable with a different value (not a copy of loop var)
for i, v := range []int{1, 2, 3} {
    doubled := v * 2
    idx := i + 1
    fmt.Println(idx, doubled)
}
```

```golang
// Without check-alias: copying to a differently-named variable is allowed by default
for i, v := range []int{1, 2, 3} {
    _i := i // Not flagged when check-alias is false (default)
    _v := v // Not flagged when check-alias is false (default)
    fmt.Println(_i, _v)
}
```

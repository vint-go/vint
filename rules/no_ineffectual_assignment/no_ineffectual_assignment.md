---
title: noIneffectualAssignment
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noIneffectualAssignment`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noIneffectualAssignment:
    check-escaping-errors: false # check escaping variables of type error, may cause false positives
```

## Details

Detects ineffectual assignments in Go code. An assignment is ineffectual if the variable assigned is not thereafter used. This can indicate dead code, logic errors, or leftover code from refactoring.

The analyzer performs control flow analysis to determine whether an assigned value is ever read before the variable is reassigned or goes out of scope. It catches cases such as:

- A variable is assigned a value that is immediately overwritten by another assignment without being read.
- A variable is assigned a value but the function returns or the variable goes out of scope before the value is ever used.
- A variable is assigned in a branch but never read in any subsequent path.

Note that this analysis operates at the syntax level and does not consider type information. As a result, assignments to struct fields are never flagged. The tool is designed to produce no false positives under normal usage.

When `check-escaping-errors` is enabled, the analyzer will also report ineffectual assignments to variables of type `error` that escape (e.g., are passed to other functions or returned). This option may produce false positives and is disabled by default.

Source: https://github.com/gordonklaus/ineffassign

## Examples

### Invalid

```golang
// Variable assigned but immediately overwritten
func example() int {
    x := 1
    x = 2
    return x
}
```

```golang
// Variable assigned but never used before function returns
func example() {
    x := computeValue()
    if condition {
        return
    }
    x = otherValue()
    fmt.Println(x)
}
```

```golang
// Variable assigned in all branches before being read
func example(flag bool) int {
    x := 0
    if flag {
        x = 1
    } else {
        x = 2
    }
    return x
}
```

```golang
// Error variable assigned but overwritten without checking
func example() error {
    err := doSomething()
    err = doSomethingElse()
    return err
}
```

```golang
// Variable assigned but shadowed by short variable declaration
func example() {
    x := 1
    fmt.Println(x)
    x = 2
    x = 3
    fmt.Println(x)
}
```

### Valid

```golang
// Variable assigned and subsequently used
func example() int {
    x := 1
    return x
}
```

```golang
// Variable assigned and read before reassignment
func example() {
    x := computeValue()
    fmt.Println(x)
    x = otherValue()
    fmt.Println(x)
}
```

```golang
// Error variable checked before reassignment
func example() error {
    err := doSomething()
    if err != nil {
        return err
    }
    err = doSomethingElse()
    return err
}
```

```golang
// Variable assigned and used in a loop
func example(items []int) int {
    sum := 0
    for _, item := range items {
        sum += item
    }
    return sum
}
```

```golang
// Blank identifier used intentionally to discard a value
func example() {
    _, err := doSomething()
    if err != nil {
        log.Fatal(err)
    }
}
```

```golang
// Variable used after if/else where branches don't touch it
func example(flag bool) int {
    x := 10
    if flag {
        fmt.Println("a")
    } else {
        fmt.Println("b")
    }
    return x
}
```

```golang
// Initial value used inside a loop
func example() int {
    count := 0
    for i := 0; i < 5; i++ {
        total := 100 + count*count
        fmt.Println(total)
        count++
    }
    return count
}
```

```golang
// Value read via append before reassignment
func example(flag bool) []int {
    items := []int{1}
    if flag {
        items = append(items, 2)
    } else {
        items = append(items, 3)
    }
    return items
}
```

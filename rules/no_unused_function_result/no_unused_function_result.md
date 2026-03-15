---
title: noUnusedFunctionResult
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnusedFunctionResult`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnusedFunctionResult:
    # rule options here
```

## Details

Checks for unused results of calls to certain pure functions. Some functions are known to have no side effects and are only useful for their return values. Calling such functions without using the result is almost certainly a mistake.

Functions checked by default include `fmt.Errorf`, `fmt.Sprintf`, `fmt.Sprint`, `errors.New`, `sort.Reverse`, and methods like `String()` and `Error()` on well-known types.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/unusedresult

## Examples

### Invalid

```golang
import "fmt"

func example() {
    // Bad: result of fmt.Sprintf is not used
    fmt.Sprintf("hello %s", "world")
}
```

```golang
import "errors"

func example() {
    // Bad: result of errors.New is not used
    errors.New("something went wrong")
}
```

```golang
import "sort"

func example() {
    s := []int{3, 1, 2}
    // Bad: result of sort.Reverse is not used
    sort.Reverse(sort.IntSlice(s))
}
```

### Valid

```golang
import "fmt"

func example() {
    // Good: result is assigned to a variable
    msg := fmt.Sprintf("hello %s", "world")
    fmt.Println(msg)
}
```

```golang
import "errors"

func example() error {
    // Good: result is returned
    return errors.New("something went wrong")
}
```

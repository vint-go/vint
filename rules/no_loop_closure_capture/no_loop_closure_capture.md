---
title: noLoopClosureCapture
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noLoopClosureCapture`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noLoopClosureCapture:
    # rule options here
```

## Details

Checks for references to enclosing loop variables from within nested functions (closures). A common mistake in Go (prior to Go 1.22) is to capture a loop variable in a goroutine or closure launched inside the loop. Because the loop variable is shared across all iterations, the closure may observe a different value than intended, typically the value from the last iteration.

Note: Starting with Go 1.22, the loop variable semantics changed so that each iteration gets its own copy. This analyzer is most relevant for code targeting Go versions before 1.22.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/loopclosure

## Examples

### Invalid

```golang
func example() {
    for _, v := range []int{1, 2, 3} {
        // Bad: v is captured by reference, all goroutines will see the last value
        go func() {
            fmt.Println(v)
        }()
    }
}
```

```golang
func example() {
    for i := 0; i < 3; i++ {
        // Bad: i is captured by reference
        defer func() {
            fmt.Println(i)
        }()
    }
}
```

### Valid

```golang
func example() {
    for _, v := range []int{1, 2, 3} {
        v := v // Good: create a new variable for each iteration
        go func() {
            fmt.Println(v)
        }()
    }
}
```

```golang
func example() {
    for _, v := range []int{1, 2, 3} {
        // Good: pass v as a parameter
        go func(val int) {
            fmt.Println(val)
        }(v)
    }
}
```

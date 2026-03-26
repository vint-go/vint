---
title: noUnreachableCode
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnreachableCode`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnreachableCode:
    # rule options here
```

## Details

Checks for unreachable code. This analyzer detects code that can never be executed because it appears after a statement that unconditionally exits the function, such as `return`, `panic`, `os.Exit`, `log.Fatal`, or an infinite loop with no break. Unreachable code is usually a sign of a programming error or dead code that should be removed.

For library function calls like `os.Exit` or `log.Fatal`, the rule verifies that the identifier actually refers to an imported package rather than a local variable with the same name. For example, a local variable named `log` (e.g., of type `*zap.SugaredLogger`) calling `.Panic()` will not be treated as the stdlib `log.Panic()`.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/unreachable

## Examples

### Invalid

```golang
func example() int {
    return 42
    fmt.Println("this is unreachable") // Bad: code after return
}
```

```golang
func example() {
    panic("fatal error")
    cleanup() // Bad: code after panic
}
```

### Valid

```golang
func example() int {
    fmt.Println("computing result")
    return 42
}
```

```golang
func example(x int) int {
    if x > 0 {
        return x
    }
    return -x // Good: reachable when x <= 0
}
```

```golang
func example(log *zap.SugaredLogger) {
    log.Panic("something went wrong")
    cleanup() // Good: log is a local variable, not the stdlib log package
}
```

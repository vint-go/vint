---
title: noExcessiveControlNesting
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/complexity/noExcessiveControlNesting`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/complexity/noExcessiveControlNesting:
    arguments:
      - 3
```

## Details

Warns if nesting level of control structures (`if-then-else`, `for`, `switch`) exceeds a given maximum.

Deeply nested control structures make code harder to read, understand, and maintain. By enforcing a maximum nesting level, developers are encouraged to refactor complex logic into smaller, more manageable functions.

The rule accepts a single integer argument specifying the maximum allowed nesting level. The default value is 5.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// With max nesting set to 2:
func example() {
    if condition1 {
        if condition2 {
            if condition3 { // error: control flow nesting exceeds 2
                doSomething()
            }
        }
    }
}
```

```golang
// With max nesting set to 2:
func example() {
    for i := range items {
        if condition {
            for { // error: control flow nesting exceeds 2
                break
            }
        }
    }
}
```

### Valid

```golang
// With max nesting set to 2:
func example() {
    if condition1 {
        if condition2 {
            doSomething()
        }
    }
}
```

```golang
// Refactored to reduce nesting:
func example() {
    if !condition1 {
        return
    }
    if condition2 {
        doSomething()
    }
}
```

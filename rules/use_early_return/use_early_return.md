---
title: useEarlyReturn
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useEarlyReturn`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useEarlyReturn:
    arguments:
      - "preserve-scope"
      - "allow-jump"
```

## Details

In Go it is idiomatic to minimize nesting statements. A typical example is to avoid if-then-else constructions where the else block deviates control flow (e.g. returns, continues, breaks). This rule spots such constructions and suggests inverting the `if` condition to reduce nesting.

Available configuration flags (passed as string arguments):

- `preserve-scope`: do not suggest refactorings that would increase variable scope.
- `allow-jump`: suggest introducing a new jump (`return`, `continue` or `break` statement) to reduce nesting. By default, only relocation of existing jumps (i.e. from the `else` clause) are suggested.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
if cond {
    // do something
} else {
    // do other thing
    return ...
}
```

```golang
// With allow-jump enabled:
func example() {
    if cond {
        println()
        println()
        println()
    }
    // Can be rewritten as: if !cond { return } ...
}
```

### Valid

```golang
if !cond {
    // do other thing
    return ...
}

// do something
```

```golang
// Already uses early return pattern
if err != nil {
    return err
}
// proceed with normal flow
```

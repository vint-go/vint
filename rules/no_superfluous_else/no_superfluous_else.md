---
title: noSuperfluousElse
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noSuperfluousElse`
- This rule is recommended, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noSuperfluousElse:
    arguments:
      - "preserveScope"
```

## Details

To improve the readability of code, it is recommended to reduce the indentation as much as possible.
This rule highlights redundant _else-blocks_ that can be eliminated from the code.

Unlike `useIndentErrorFlow` which handles `if` blocks ending with a `return` statement, this rule handles `if` blocks ending with other control flow deviations such as `continue`, `break`, `goto`, `panic`, `log.Fatal`, and `os.Exit`.

Configuration: ([]string) rule flags. Available flags are:

- `preserveScope`: do not suggest refactorings that would increase variable scope

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
for {
    if f() {
        continue
    } else { // superfluous else block
        log.Printf("non-positive")
    }
}
```

```golang
if f() {
    log.Fatal("x")
} else { // superfluous else block
    log.Printf("non-positive")
}
```

### Valid

```golang
for {
    if f() {
        continue
    }
    log.Printf("non-positive")
}
```

```golang
if f() {
    log.Fatal("x")
}
log.Printf("non-positive")
```

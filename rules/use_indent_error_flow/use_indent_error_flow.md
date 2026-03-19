---
title: useIndentErrorFlow
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useIndentErrorFlow`
- This rule is recommended, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useIndentErrorFlow:
    arguments:
      - "preserve-scope"
```

## Details

To improve the readability of code, it is recommended to reduce the indentation as much as possible.
This rule highlights redundant _else-blocks_ that can be eliminated from the code when the if-block ends with a return statement.

More information: [Go Code Review Comments - Indent Error Flow](https://go.dev/wiki/CodeReviewComments#indent-error-flow).

Available configuration flags (passed as string arguments):

- `preserve-scope`: do not suggest refactorings that would increase variable scope.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
if x > 0 {
    return true
} else {
    log.Printf("non-positive x: %d", x)
}
```

```golang
if ok := f(); ok {
    return "it's okay"
} else {
    return "it's NOT okay!"
}
```

### Valid

```golang
if x > 0 {
    return true
}
log.Printf("non-positive x: %d", x)
```

```golang
ok := f()
if ok {
    return "it's okay"
}
return "it's NOT okay!"
```

---
title: noDuplicateBranchBody
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noDuplicateBranchBody`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noDuplicateBranchBody:
    # no additional options
```

## Details

Detects duplicated branch bodies inside conditional statements. This checker identifies `if` statements where both the then-branch and else-branch contain identical code bodies. When both branches do the same thing, the conditional is redundant and likely indicates a logic error.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
if condition {
    fmt.Println("result")
} else {
    fmt.Println("result")
}
```

### Valid

```golang
if condition {
    fmt.Println("true result")
} else {
    fmt.Println("false result")
}
```

```golang
// If both branches truly need the same code, remove the conditional
fmt.Println("result")
```

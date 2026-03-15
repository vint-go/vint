---
title: noDuplicateCase
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noDuplicateCase`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noDuplicateCase:
    # no additional options
```

## Details

Detects duplicated case clauses inside switch or select statements. When the same expression appears in multiple case clauses of a switch statement, or duplicate communication operations appear in a select statement, one of the cases is likely a mistake.

This rule is distinct from `lint/suspicious/noRedundantBoolCondition`, which detects duplicate conditions within `&&` and `||` boolean expressions. This rule specifically targets duplicate case values in switch and select statements.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
switch x {
case 1:
    handleOne()
case 2:
    handleTwo()
case 1: // duplicate of the first case
    handleOneAgain()
}
```

### Valid

```golang
switch x {
case 1:
    handleOne()
case 2:
    handleTwo()
case 3:
    handleThree()
}
```

---
title: noSideEffectInInitClause
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noSideEffectInInitClause`
- This rule is not recommended (experimental, opinionated).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noSideEffectInInitClause:
    # no additional options
```

## Details

Detects non-assignment statements inside if/switch init clause. This checker identifies problematic code patterns where non-assignment statements (such as side-effect function calls) appear in the initialization clause of if or switch statements. It recommends moving the statement outside the conditional, before the if/switch block, for improved code clarity.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
if sideEffect(); condition {
    doSomething()
}
```

### Valid

```golang
sideEffect()
if condition {
    doSomething()
}
```

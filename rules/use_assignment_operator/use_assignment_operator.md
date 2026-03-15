---
title: useAssignmentOperator
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useAssignmentOperator`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useAssignmentOperator:
    # no additional options
```

## Details

Detects assignments that can be simplified by using assignment operators. When a variable is used on both sides of an assignment with a binary operator, the expression can be written more concisely using the corresponding assignment operator (`+=`, `-=`, `*=`, etc.).

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
x = x + 1
```

```golang
y = y * 2
```

```golang
count = count - delta
```

### Valid

```golang
x += 1
```

```golang
y *= 2
```

```golang
count -= delta
```

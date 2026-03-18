---
title: noUnnecessaryBlock
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noUnnecessaryBlock`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noUnnecessaryBlock:
    # no additional options
```

## Details

Detects unnecessary braced statement blocks. This checker identifies two patterns:

1. Block statements within statement lists that contain no variable definitions or declarations -- the braces can be removed.
2. Case/switch statements with a single block statement -- the block is redundant in a case clause.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
func f() {
    {
        doSomething()
        doSomethingElse()
    }
}
```

```golang
switch x {
case 1:
    {
        handleOne()
    }
}
```

### Valid

```golang
func f() {
    doSomething()
    doSomethingElse()
}
```

```golang
func f() {
    {
        x := computeValue() // block is needed for scoping
        use(x)
    }
}
```

```golang
switch x {
case 1:
    handleOne()
}
```

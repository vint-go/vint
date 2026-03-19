---
title: noUselessFallthrough
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noUselessFallthrough`
- This rule is **not** enabled by default. Enable it by setting it in the configuration.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noUselessFallthrough:
    # no configuration options
```

## Details

This rule warns on useless `fallthrough` statements in case clauses of switch statements. A `fallthrough` is considered useless if it is the single statement of a case clause body.

Go allows `switch` statements with clauses that group multiple cases. Therefore, a case clause whose only statement is `fallthrough` to continue to the next case can be simplified by merging the case expressions into a single clause.

When a `fallthrough` statement has a comment attached, the rule still reports it but with lower confidence (0.5 instead of 1.0). Fallthroughs to a `default` case are not flagged.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
switch category {
case "Lu":
    fallthrough // useless: can consolidate with next case
case "Ll":
    fallthrough // useless: can consolidate with next case
case "Lt":
    return true
default:
    return false
}
```

The above can be simplified to:

```golang
switch category {
case "Lu", "Ll", "Lt":
    return true
default:
    return false
}
```

```golang
switch a {
case 0:
    fallthrough // useless fallthrough
case 1:
    println()
}
```

### Valid

```golang
switch a {
case 0:
    println()
    fallthrough // not useless: body has more than one statement
default:
}
```

```golang
switch a {
case 0:
    fallthrough // valid: falls through to default
default:
    println()
}
```

```golang
switch a {
case 0:
    if condition {
        fallthrough
    }
    doSomething()
case 1:
    println()
}
```

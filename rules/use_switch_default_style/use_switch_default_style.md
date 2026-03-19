---
title: useSwitchDefaultStyle
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useSwitchDefaultStyle`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useSwitchDefaultStyle:
    arguments:
      - "allowNoDefault"
      - "allowDefaultNotLast"
```

## Details

This rule enforces consistent usage of `default` on `switch` statements. It can check for `default` case clause occurrence and/or position in the list of case clauses.

Configuration options (non-mutually exclusive):

- `allowNoDefault`: allows `switch` without `default` case clause.
- `allowDefaultNotLast`: allows `default` case clause to be not the last clause of the `switch`.

By default (no options), the rule enforces that all `switch` statements have a `default` clause as its last case clause.

Using `allowDefaultNotLast` alone enforces that all `switch` statements have a `default` clause but its position is unimportant.

Using `allowNoDefault` alone enforces that in all `switch` statements with a `default` clause, the `default` is the last case clause.

Notice that a configuration including both options will effectively deactivate the whole rule.

If all case branches end with a jump statement (`return` or `break`), the rule does not require a `default` clause.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// Missing default case clause
switch expression {
case condition:
}
```

```golang
// Default case clause is not the last one
switch expression {
default:
case condition:
}
```

### Valid

```golang
// Default case clause is last
switch expression {
case condition:
default:
}
```

```golang
// All branches end with jump statements (no default required)
switch expression {
case condition1:
    return
case condition2:
    break
}
```

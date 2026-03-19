---
title: noUselessBreak
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noUselessBreak`
- This rule is **not** enabled by default. Enable it by setting it in the configuration.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noUselessBreak:
    # no configuration options
```

## Details

This rule warns on useless `break` statements in case clauses of switch and select statements. Go, unlike other programming languages like C, only executes statements of the selected case while ignoring the subsequent case clauses. Therefore, inserting a `break` at the end of a case clause has no effect.

Because `break` statements are rarely used in case clauses, when switch or select statements are inside a for-loop, the programmer might wrongly assume that a `break` in a case clause will take the control out of the loop. The rule emits a specific warning for such cases.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
switch {
case true:
    break // useless break in case clause
}
```

```golang
for {
    switch {
    case c1:
        break // useless break - affects the switch, not the enclosing loop
    }
}
```

```golang
select {
case c:
    break // useless break in case clause
}
```

### Valid

```golang
switch {
case true:
    if someCondition {
        break // not the last statement, conditional break is fine
    }
    doSomething()
}
```

```golang
switch val.Kind() {
case reflect.Array:
    if val.Len() == 0 {
        break // break used conditionally in the middle of a case
    }
    process(val)
    return
}
```

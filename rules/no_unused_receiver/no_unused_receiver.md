---
title: noUnusedReceiver
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noUnusedReceiver`
- This rule is not recommended. Enable it explicitly via configuration.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noUnusedReceiver:
    arguments:
      - allowRegex: "^_" # Optional regex for allowed unused receiver name patterns. Default: "^_$" (only blank identifier).
```

## Details

This rule warns on unused method receivers. Methods with unused receivers can be a symptom of an unfinished refactoring or a bug. If a method does not reference its receiver, it may indicate that the method should be a plain function instead, or that the receiver was accidentally left unused after a code change.

By default, receivers named `_` are allowed. You can configure an `allowRegex` option to specify additional allowed unused receiver name patterns. The option name can be written as `allowRegex`, `allowregex`, or `allow-regex`.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// The receiver 'f' is never used in the method body.
func (f *Formatter) Name() string {
    return "formatter"
}
```

```golang
// The receiver 'u' is shadowed inside the loop and never
// actually referenced as the method receiver.
func (u *Unix) Foo() (string, error) {
    for item := range items {
        u := 1 // shadows the receiver
        fmt.Printf("%v\n", u)
    }
    return "", nil
}
```

### Valid

```golang
// The receiver is explicitly named '_', indicating it is intentionally unused.
func (_ *Formatter) Name() string {
    return "formatter"
}
```

```golang
// The receiver 'u' is referenced in the return statement.
func (u *Unix) Foos() (string, error) {
    for item := range items {
        u := 1
        fmt.Printf("%v\n", u)
    }
    return u, nil
}
```

```golang
// The receiver 'u' is used to modify a field.
func (u *Unix) Bar() (string, error) {
    for item := range items {
        u.path = nil
        fmt.Printf("%v\n", u)
    }
    return "", nil
}
```

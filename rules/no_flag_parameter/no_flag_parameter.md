---
title: noFlagParameter
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noFlagParameter`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noFlagParameter:
    # no configuration options
```

## Details

If a function controls the flow of another by passing it information on what to do, both functions are said to be [control-coupled](https://en.wikipedia.org/wiki/Coupling_(computer_programming)#Procedural_programming).
Coupling among functions must be minimized for better maintainability of the code.
This rule warns on boolean parameters that create a control coupling.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
func foo(a bool, b int) {
	if a {
		// do something
	}
}
```

### Valid

```golang
func bar(a bool, b int) {
	str := mystruct{a, b}
}
```

```golang
func baz(a int, b bool) {
	lBool := true
	if lBool {
		// do something
	}
}
```

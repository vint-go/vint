---
title: noModifiedParameter
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noModifiedParameter`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noModifiedParameter:
    # rule options here
```

## Details

A function that modifies its parameters can be hard to understand. It can also be misleading if the arguments are passed by value by the caller. This rule warns when a function modifies one or more of its parameters or when parameters are passed to functions that modify them (e.g. `slices.Delete`).

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
func process(a int) {
	a = 42 // modifies parameter 'a'
}
```

```golang
func increment(x int) {
	x++ // modifies parameter 'x'
}
```

```golang
func removeFirst(s []int) {
	s = slices.Delete(s, 0, 1) // modifies parameter 's' via slices.Delete
}
```

### Valid

```golang
func process(a int) {
	b := a + 1 // uses parameter without modifying it
	fmt.Println(b)
}
```

```golang
func removeFirst(s []int) []int {
	s2 := slices.Clone(s)
	s2 = slices.Delete(s2, 0, 1) // modifies a copy, not the parameter
	return s2
}
```

```golang
func update(s *foo) {
	s.a = "value" // modifying a field of a pointer parameter is fine
}
```

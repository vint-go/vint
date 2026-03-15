---
title: noSelfAssignment
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSelfAssignment`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSelfAssignment:
    # rule options here
```

## Details

Detects useless assignments. This analyzer checks for self-assignments of the form `x = x`, which are almost always a mistake. The programmer likely intended to assign to a different variable or use a different value on the right-hand side.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/assign

## Examples

### Invalid

```golang
func example() {
    x := 5
    x = x // useless self-assignment
}
```

```golang
func example(s *MyStruct) {
    s.Field = s.Field // useless self-assignment of struct field
}
```

### Valid

```golang
func example() {
    x := 5
    y := x // assignment to a different variable
}
```

```golang
func example(s *MyStruct) {
    s.Field = computeNewValue() // assignment with a different value
}
```

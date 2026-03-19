---
title: noUnnecessaryStmt
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noUnnecessaryStmt`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noUnnecessaryStmt:
    # rule options here
```

## Details

This rule suggests removing redundant statements like a `break` at the end of a case block or a bare `return` at the end of a function with no return values, for improving the code's readability. It also detects `switch` statements with only one case that can be replaced by an `if-then`.

In Go, `case` blocks in `switch` statements do not fall through by default, so a trailing `break` is unnecessary. Similarly, functions without return values do not need an explicit `return` at the end.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
func bar() {
	a := 1
	switch a {
	case 1:
		a++
		break // unnecessary break at the end of case clause
	}

	return // unnecessary return statement in function with no return values
}
```

```golang
func baz(x int) {
	switch x { // switch with only one case can be replaced by an if-then
	case 1:
		println("one")
	}
}
```

### Valid

```golang
func bar() {
	a := 1
	switch a {
	case 1:
		a++
	case 2:
		println("two")
	}
}
```

```golang
func baz(x int) {
	if x == 1 {
		println("one")
	}
}
```

```golang
func qux() int {
	return 42 // return with value is fine
}
```

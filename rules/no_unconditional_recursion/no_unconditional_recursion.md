---
title: noUnconditionalRecursion
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnconditionalRecursion`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnconditionalRecursion:
    # no configuration options
```

## Details

Unconditional recursive calls will produce infinite recursion, thus program stack overflow.
This rule detects and warns about unconditional (direct) recursive calls.

The rule walks function bodies looking for calls to the function itself that are not guarded
by any conditional control structure (if, for with condition, switch, select, range).
It also checks whether conditional branches contain control-flow exits (return, panic,
os.Exit, log.Fatal, etc.) that would prevent the recursion from being truly unconditional.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
func foo() {
	foo() // unconditional recursive call
}
```

```golang
func (mt *myType) Foo() int {
	return mt.Foo() // unconditional recursive call on method
}
```

```golang
func bar() {
	for {
		bar() // recursive call inside unconditional loop
	}
}
```

```golang
func baz() {
	go baz() // produces infinite number of goroutines
}
```

### Valid

```golang
func foo(n int) {
	if n <= 0 {
		return
	}
	foo(n - 1) // conditional: guarded by if-return
}
```

```golang
func bar() {
	if true {
		panic("done")
	}
	bar() // conditional exit seen before recursive call
}
```

```golang
func baz() {
	_ = callback(func() {
		baz() // inside a closure, not necessarily unconditional
	})
}
```

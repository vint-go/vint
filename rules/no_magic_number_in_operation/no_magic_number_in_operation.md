---
title: noMagicNumberInOperation
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noMagicNumberInOperation`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noMagicNumberInOperation:
    ignored-numbers: "0,0.0,1,1.0"
    ignored-files: "_test.go"
```

## Details

Detects magic numbers used in arithmetic and binary operations. A magic number is a numeric literal that is not defined as a constant, but which may change, and therefore can be hard to update. Using magic numbers directly in operations makes code less readable and harder to maintain, as the purpose or meaning of the numeric literal in the computation is not clear.

This check examines numeric literals appearing in binary expressions (such as multiplication, addition, subtraction, and division) within both assignment statements and parenthesized expressions. It inspects both operands of binary expressions and also flags magic numbers in nested or chained operations. Multiple magic numbers in a single expression are each reported individually.

By default, the numbers `0`, `0.0`, `1`, and `1.0` are excluded from detection. Test files (`_test.go`) are also excluded by default.

Source: https://github.com/tommy-muehle/go-mnd

## Examples

### Invalid

```golang
package example

func compute(y int) {
	_ = y * 20 // Magic number: 20, in <operation> detected
}
```

```golang
package example

func compute2(y int) {
	_ = 10 * y // Magic number: 10, in <operation> detected
}
```

```golang
package example

func compute3(y int) {
	_ = 5 * y * 6 // Magic number: 5, in <operation> detected
	               // Magic number: 6, in <operation> detected
}
```

```golang
package example

func compute4() {
	const c = 24
	_ = c + (42 * 10) // Magic number: 42, in <operation> detected
	                   // Magic number: 10, in <operation> detected
}
```

```golang
package example

func compute5(x int32) {
	if (42 * x) > 10 { // Magic number: 42, in <operation> detected
		               // Magic number: 10, in <operation> detected
	}
}
```

```golang
package example

func compute6(x float32) {
	if 10 < (x * 1.0) { // Magic number: 10, in <operation> detected
		                 // 1.0 is excluded by default
	}
}
```

### Valid

```golang
package example

const multiplier = 20

func compute(y int) {
	_ = y * multiplier
}
```

```golang
package example

const (
	factor = 42
	base   = 10
)

func compute2() {
	const c = 24
	_ = c + (factor * base)
}
```

```golang
package example

func computeSimple(y int) {
	_ = y * 1 // 1 is excluded by default
	_ = y + 0 // 0 is excluded by default
}
```

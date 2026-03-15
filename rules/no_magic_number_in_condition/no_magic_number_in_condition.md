---
title: noMagicNumberInCondition
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noMagicNumberInCondition`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noMagicNumberInCondition:
    ignored-numbers: "0,0.0,1,1.0"
    ignored-files: "_test.go"
```

## Details

Detects magic numbers used in conditional expressions (`if` statements). A magic number is a numeric literal that is not defined as a constant, but which may change, and therefore can be hard to update. Using magic numbers directly in conditions makes code less readable and harder to maintain, as the threshold or boundary value has no descriptive name.

This check examines both the left and right operands of binary comparison expressions within `if` statement conditions. It flags any numeric literal (integer or floating-point) that is not in the ignored numbers list.

By default, the numbers `0`, `0.0`, `1`, and `1.0` are excluded from detection. Test files (`_test.go`) are also excluded by default.

Source: https://github.com/tommy-muehle/go-mnd

## Examples

### Invalid

```golang
package example

func checkValue(x int) {
	if x > 7 { // Magic number: 7, in <condition> detected
		// do something
	}
}
```

```golang
package example

func checkValue2(x int) {
	if 8 > x { // Magic number: 8, in <condition> detected
		// do something
	}
}
```

### Valid

```golang
package example

const maxRetries = 7

func checkValue(x int) {
	if x > maxRetries {
		// do something
	}
}
```

```golang
package example

func checkFloat(x float32) {
	if x > 1.0 { // 1.0 is excluded by default
	}

	if x < 0.0 { // 0.0 is excluded by default
	}
}
```

```golang
package example

func checkString(y string) {
	if "test" == y { // string comparisons are not checked
	}
}
```

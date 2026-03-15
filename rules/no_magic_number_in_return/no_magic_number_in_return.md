---
title: noMagicNumberInReturn
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noMagicNumberInReturn`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noMagicNumberInReturn:
    ignored-numbers: "0,0.0,1,1.0"
    ignored-files: "_test.go"
```

## Details

Detects magic numbers used in `return` statements. A magic number is a numeric literal that is not defined as a constant, but which may change, and therefore can be hard to update. Using magic numbers directly in return statements makes code less readable and harder to maintain, as the meaning of the returned value is not clear without a descriptive constant name.

This check examines both direct numeric literals in return statements and numeric literals within binary expressions used as return values. String literals that happen to contain digits (such as `"3"`) are not flagged. Multiple magic numbers in a single return expression are each reported individually.

By default, the numbers `0`, `0.0`, `1`, and `1.0` are excluded from detection. Test files (`_test.go`) are also excluded by default.

Source: https://github.com/tommy-muehle/go-mnd

## Examples

### Invalid

```golang
package example

func getCount() int {
	return 3 // Magic number: 3, in <return> detected
}
```

```golang
package example

func getPi() float64 {
	return 2.0 // Magic number: 2.0, in <return> detected
}
```

```golang
package example

func addOffset(x int32) int32 {
	return x + 42 // Magic number: 42, in <return> detected
}
```

```golang
package example

func addOffset2(x int32) int32 {
	return 42 + x // Magic number: 42, in <return> detected
}
```

```golang
package example

func complexReturn(x int32) int32 {
	return x + (42 * 10) // Magic number: 42, in <return> detected
	                     // Magic number: 10, in <return> detected
}
```

```golang
package example

func complexReturn2(x int32) int32 {
	return (42 * x) + 10 // Magic number: 42, in <return> detected
	                     // Magic number: 10, in <return> detected
}
```

```golang
package example

func getFloat() float32 {
	return 3.0 // Magic number: 3.0, in <return> detected
}
```

### Valid

```golang
package example

const defaultCount = 3

func getCount() int {
	return defaultCount
}
```

```golang
package example

const offset = 42

func addOffset(x int32) int32 {
	return x + offset
}
```

```golang
package example

func getString() string {
	return "3" // string literals are not checked
}
```

```golang
package example

func getZero() float32 {
	return 0.0 // 0.0 is excluded by default
}
```

```golang
package example

func getOne() float32 {
	return 1.0 // 1.0 is excluded by default
}
```

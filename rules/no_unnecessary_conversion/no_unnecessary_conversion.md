---
title: noUnnecessaryConversion
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noUnnecessaryConversion`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a **fix**.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noUnnecessaryConversion:
    fast-math: false
    safe: false
```

### Options

#### `fast-math`
Default: `false`

When set to `true`, the linter will also flag explicit floating-point and complex number type conversions as unnecessary. By default (`false`), these conversions are preserved because in Go (since Go 1.9) explicit floating-point conversions are always significant: they force rounding and prevent operation fusing by the compiler. Enabling `fast-math` means you accept that removing such conversions may subtly change floating-point arithmetic results.

#### `safe`
Default: `false`

When set to `true`, the linter performs additional context-aware analysis before flagging a conversion as unnecessary. It examines the surrounding context (assignments, binary expressions, function call arguments, return statements, shift operations, composite literals, and pointer/unary expressions) to ensure that removing the conversion would not change the semantics of the code. This is a more conservative mode that reduces false positives.

## Details

The `noUnnecessaryConversion` rule detects unnecessary type conversions in Go code. A type conversion expression `T(x)` is considered unnecessary when `x` already has type `T`, meaning the conversion is a no-op that adds visual noise without changing the program's behavior.

The linter analyzes the code by walking the AST (Abstract Syntax Tree) and checking each type conversion call expression. For a conversion to be flagged as unnecessary, the following conditions must hold:

1. The expression must be a call expression with exactly one argument (a type conversion, not a regular function call).
2. No ellipsis (`...`) is present in the call.
3. The callee must be a type expression, not a function.
4. The argument type must be identical to the target conversion type (using `types.Identical`).
5. The argument must not be an untyped value (due to a Go compiler workaround for [golang.org/issue/13061](https://github.com/golang/go/issues/13061)).
6. Unless `fast-math` is enabled, floating-point and complex number conversions (`float32`, `float64`, `complex64`, `complex128`) are not flagged because they force rounding and prevent operation fusing.
7. When `safe` mode is enabled, additional context checks verify that removing the conversion would be semantically safe in the surrounding code.

Source: https://github.com/mdempsky/unconvert

## Examples

### Invalid

```golang
// Unnecessary conversion of an int variable to int
var x int = 42
y := int(x) // unnecessary conversion
```

```golang
// Unnecessary conversion of a string variable to string
var name string = "hello"
s := string(name) // unnecessary conversion
```

```golang
// Unnecessary conversion of a bool variable to bool
var flag bool = true
b := bool(flag) // unnecessary conversion
```

```golang
// Unnecessary conversion of a custom type to itself
type Counter int64
var c Counter = 10
d := Counter(c) // unnecessary conversion
```

```golang
// Unnecessary pointer conversion
var p *int
q := (*int)(p) // unnecessary conversion
```

```golang
// Unnecessary conversion in an interface type
var w io.Writer
x := io.Writer(w) // unnecessary conversion
```

```golang
// Unnecessary conversion of a function type
type Constructor func() ID
var ctor Constructor
f := Constructor(ctor) // unnecessary conversion
```

### Valid

```golang
// Converting between different types is necessary
var x int32 = 42
y := int64(x) // different types, conversion is needed
```

```golang
// Converting an untyped constant to a specific type
const maxSize = 100
y := int64(maxSize) // untyped constant needs explicit conversion
```

```golang
// Float-to-float conversion preserves rounding semantics (without fast-math)
var f1 float64 = 1.1
var f2 float64 = 2.2
result := float64(f1 * f2) // forces rounding, prevents operation fusing
```

```golang
// Converting between a named type and its underlying type
type UserID string
var name string = "alice"
id := UserID(name) // different types, conversion is needed
```

```golang
// Converting from one numeric type to another
var count uint32 = 5
index := int(count) // different types, conversion is needed
```

```golang
// Complex number conversion preserves precision semantics (without fast-math)
var c1 complex128 = complex(1.0, 2.0)
c2 := complex128(c1) // preserves precision semantics
```

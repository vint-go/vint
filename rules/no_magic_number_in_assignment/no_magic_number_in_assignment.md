---
title: noMagicNumberInAssignment
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noMagicNumberInAssignment`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noMagicNumberInAssignment:
    ignored-numbers: "0,0.0,1,1.0"
    ignored-files: "_test.go"
```

## Details

Detects magic numbers used in assignment statements and struct literal field values. A magic number is a numeric literal that is not defined as a constant, but which may change, and therefore can be hard to update. Using magic numbers directly in assignments makes code less readable and harder to maintain, as the purpose or meaning of the number is not immediately clear.

This check covers direct variable assignments, struct field initializations in composite literals, key-value expressions, unary expressions (such as negative numbers), and binary expressions on the right-hand side of assignments.

By default, the numbers `0`, `0.0`, `1`, and `1.0` are excluded from detection. Test files (`_test.go`) are also excluded by default.

Source: https://github.com/tommy-muehle/go-mnd

## Examples

### Invalid

```golang
package example

func createOrder() interface{} {
	s := struct {
		Amount int
	}{
		Amount: 100, // Magic number: 100, in <assign> detected
	}
	return s
}
```

```golang
package example

import (
	"net/http"
	"time"
)

func createClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second, // Magic number: 5, in <assign> detected
	}
}
```

```golang
package example

import "fmt"

func compute() {
	res := -12 // Magic number: 12, in <assign> detected
	fmt.Println(res)
}
```

```golang
package example

import "fmt"

func compute2() {
	var x int32
	res := x + -12 // Magic number: 12, in <assign> detected
	fmt.Println(res)
}
```

```golang
package example

import "fmt"

func compute3() {
	var x int32
	res := 12 + x // Magic number: 12, in <assign> detected
	fmt.Println(res)
}
```

### Valid

```golang
package example

const defaultAmount = 100

func createOrder() interface{} {
	s := struct {
		Amount int
	}{
		Amount: defaultAmount,
	}
	return s
}
```

```golang
package example

import (
	"net/http"
	"time"
)

const clientTimeout = 5

func createClient() *http.Client {
	return &http.Client{
		Timeout: clientTimeout * time.Second,
	}
}
```

```golang
package example

import "fmt"

const offset = 12

func compute() {
	var x int32
	res := x + -offset
	fmt.Println(res)
}
```

```golang
package example

func defaultValues() {
	x := 0   // 0 is excluded by default
	y := 1.0 // 1.0 is excluded by default
	_ = x
	_ = y
}
```

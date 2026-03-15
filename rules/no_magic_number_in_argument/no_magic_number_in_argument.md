---
title: noMagicNumberInArgument
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noMagicNumberInArgument`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noMagicNumberInArgument:
    ignored-numbers: "0,0.0,1,1.0"
    ignored-functions: "math.*,http.StatusText,strconv.ParseInt,strconv.ParseUint,strconv.ParseFloat,strconv.FormatInt,strconv.FormatUint,strconv.FormatFloat,time.Date"
    ignored-files: "_test.go"
```

## Details

Detects magic numbers used as function or method call arguments. A magic number is a numeric literal that is not defined as a constant, but which may change, and therefore can be hard to update. Using magic numbers directly in function arguments makes code less readable and harder to maintain, as the meaning of the number is not immediately clear from context.

By default, the numbers `0`, `0.0`, `1`, and `1.0` are excluded from detection. Test files (`_test.go`) are also excluded by default. Additionally, certain standard library functions such as `time.Date`, `strconv.ParseInt`, `strconv.ParseUint`, `strconv.ParseFloat`, `strconv.FormatInt`, `strconv.FormatUint`, and `strconv.FormatFloat` are ignored by default because they commonly accept numeric literals whose meaning is clear from the function context.

Source: https://github.com/tommy-muehle/go-mnd

## Examples

### Invalid

```golang
package example

import "math"

func calculateArea() float64 {
	return math.Abs(9.5) // Magic number: 9.5, in <argument> detected
}
```

```golang
package example

import "net/http"

func getStatusText() string {
	return http.StatusText(200) // Magic number: 200, in <argument> detected
}
```

```golang
package example

import (
	"context"
	"time"
)

func doWork() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Magic number: 5, in <argument> detected
	defer cancel()
	_ = ctx
}
```

```golang
package example

func process() {
	foobar(0, 3) // Magic number: 3, in <argument> detected
}

func foobar(a, b int) {}
```

```golang
package example

import "sync"

func createPool() {
	ch := make(chan int, 500) // Magic number: 500, in <argument> detected
	_ = ch
}
```

### Valid

```golang
package example

import "math"

const threshold = 9.5

func calculateArea() float64 {
	return math.Abs(threshold)
}
```

```golang
package example

import "net/http"

const statusOK = 200

func getStatusText() string {
	return http.StatusText(statusOK)
}
```

```golang
package example

import "os"

func shutdown() {
	os.Exit(1) // 1 is excluded by default
}
```

```golang
package example

func process() {
	foobar(0, 1) // 0 and 1 are excluded by default
}

func foobar(a, b int) {}
```

```golang
package example

import "time"

func createDate() time.Time {
	return time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC) // time.Date is ignored by default
}
```

---
title: noTypeAssertElseMisread
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noTypeAssertElseMisread`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noTypeAssertElseMisread:
    # rule options here
```

## Details

`else` branch of a type assertion is probably not reading the right value.

In a two-value type assertion `v, ok := x.(T)`, if `ok` is false, `v` will be the zero value of type `T`, not the original value `x`. Using `v` in the `else` branch may not give you the value you expect.

Source: https://staticcheck.dev/docs/checks/#SA9008

## Examples

### Invalid

```golang
package main

import "fmt"

func process(x interface{}) {
    v, ok := x.(string)
    if ok {
        fmt.Println("string:", v)
    } else {
        // v is "" (zero value), not the original x
        fmt.Println("not a string:", v)
    }
}
```

### Valid

```golang
package main

import "fmt"

func process(x interface{}) {
    v, ok := x.(string)
    if ok {
        fmt.Println("string:", v)
    } else {
        // Use x, not v
        fmt.Println("not a string:", x)
    }
}
```

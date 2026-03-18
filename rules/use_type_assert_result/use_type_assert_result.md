---
title: useTypeAssertResult
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useTypeAssertResult`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useTypeAssertResult:
    # rule options here
```

## Details

Use result of type assertion to simplify cases.

In a type switch, the switched variable is already typed within each case clause. There is no need to perform a type assertion again.

Source: https://staticcheck.dev/docs/checks/#S1034

## Examples

### Invalid

```golang
package main

import "fmt"

func process(x interface{}) {
    switch x.(type) {
    case int:
        fmt.Println(x.(int))
    case string:
        fmt.Println(x.(string))
    }
}
```

### Valid

```golang
package main

import "fmt"

func process(x interface{}) {
    switch v := x.(type) {
    case int:
        fmt.Println(v)
    case string:
        fmt.Println(v)
    }
}
```

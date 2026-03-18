---
title: noUnmarshalableType
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnmarshalableType`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnmarshalableType:
    # rule options here
```

## Details

Cannot marshal channels or functions.

`encoding/json`, `encoding/xml`, and similar packages cannot marshal channels or functions. Attempting to do so will result in a runtime error. These types have no meaningful serialized representation.

Source: https://staticcheck.dev/docs/checks/#SA1026

## Examples

### Invalid

```golang
package main

import "encoding/json"

func main() {
    ch := make(chan int)
    // Cannot marshal a channel
    data, _ := json.Marshal(ch)
    _ = data
}
```

### Valid

```golang
package main

import "encoding/json"

type Data struct {
    Name  string `json:"name"`
    Value int    `json:"value"`
}

func main() {
    d := Data{Name: "test", Value: 42}
    data, _ := json.Marshal(d)
    _ = data
}
```

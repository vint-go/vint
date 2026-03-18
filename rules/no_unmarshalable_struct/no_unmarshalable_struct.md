---
title: noUnmarshalableStruct
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noUnmarshalableStruct`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noUnmarshalableStruct:
    # rule options here
```

## Details

Trying to marshal a struct with no exported fields or misusing struct tags.

When marshaling a struct to JSON, XML, or similar formats, only exported fields are considered. A struct with no exported fields will marshal to an empty object, which is usually not the intended behavior.

Source: https://staticcheck.dev/docs/checks/#SA9005

## Examples

### Invalid

```golang
package main

import "encoding/json"

type config struct {
    name  string `json:"name"`  // unexported field, ignored
    value int    `json:"value"` // unexported field, ignored
}

func main() {
    c := config{name: "test", value: 42}
    data, _ := json.Marshal(c)
    _ = data // produces "{}"
}
```

### Valid

```golang
package main

import "encoding/json"

type Config struct {
    Name  string `json:"name"`
    Value int    `json:"value"`
}

func main() {
    c := Config{Name: "test", Value: 42}
    data, _ := json.Marshal(c)
    _ = data // produces {"name":"test","value":42}
}
```

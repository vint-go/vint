---
title: noNonPointerUnmarshal
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNonPointerUnmarshal`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNonPointerUnmarshal:
    # rule options here
```

## Details

Checks for passing non-pointer or non-interface types to unmarshal and decode functions. Functions like `json.Unmarshal`, `xml.Unmarshal`, `json.Decoder.Decode`, and similar functions require a pointer argument so they can modify the value. Passing a non-pointer will either cause a compile error or a runtime error, depending on the function.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/unmarshal

## Examples

### Invalid

```golang
import "encoding/json"

type User struct {
    Name string `json:"name"`
}

func example() {
    var u User
    data := []byte(`{"name":"Alice"}`)
    // Bad: passing non-pointer to Unmarshal
    json.Unmarshal(data, u)
}
```

### Valid

```golang
import "encoding/json"

type User struct {
    Name string `json:"name"`
}

func example() {
    var u User
    data := []byte(`{"name":"Alice"}`)
    // Good: passing a pointer to Unmarshal
    json.Unmarshal(data, &u)
}
```

```golang
import "encoding/json"

func example() {
    data := []byte(`{"name":"Alice"}`)
    // Good: passing an interface value (map is a reference type)
    var m map[string]interface{}
    json.Unmarshal(data, &m)
}
```

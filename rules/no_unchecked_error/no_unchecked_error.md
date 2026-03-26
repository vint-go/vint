---
title: noUncheckedError
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUncheckedError`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUncheckedError:
    check-blank: false
    disable-default-exclusions: false
    exclude-functions: []
```

## Details

Detects when error return values from function calls are silently ignored. In Go, many functions return an `error` value to indicate failure. Silently discarding these errors can mask bugs, data corruption, or resource leaks, making programs unreliable and difficult to debug.

This is the primary check performed by errcheck. It flags any function call whose error return value is not assigned to a variable. The rule does not perform deeper analysis on how the assigned error is subsequently handled -- it only verifies that the error is acknowledged.

By default, assigning an error to the blank identifier (`_`) is **not** flagged. To also flag blank-identifier assignments (e.g., `f, _ := os.Open("file.txt")`), set `check-blank: true`. This matches errcheck's default behavior, where `check-blank` defaults to `false`.

By default, errcheck excludes a set of standard library functions whose error returns are commonly considered non-critical, such as:

- `bytes.Buffer` and `strings.Builder` write methods (`Write`, `WriteByte`, `WriteRune`, `WriteString`)
- `fmt.Print`, `fmt.Printf`, `fmt.Println` and `fmt.Fprint` variants targeting buffers, `strings.Builder`, and `os.Stderr`
- `math/rand.Read` and `crypto/rand.Read`
- `hash.Hash.Write` and `hash/maphash.Hash.Write`
- `io.PipeReader.CloseWithError` and `io.PipeWriter.CloseWithError`

These default exclusions can be disabled via the `disable-default-exclusions` option. Additional functions can be excluded using the `exclude-functions` list.

Source: https://github.com/kisielk/errcheck

## Examples

### Invalid

```golang
// Error return value from os.Open is silently ignored
package main

import "os"

func main() {
    os.Open("file.txt")
}
```

```golang
// Error return value from file.Close is silently ignored
package main

import "os"

func main() {
    f, _ := os.Open("file.txt")
    defer f.Close()
}
```

```golang
// With check-blank: true, error assigned to blank identifier is also flagged
package main

import "os"

func main() {
    f, _ := os.Open("file.txt") // flagged only when check-blank is true
    _ = f
}
```

```golang
// Error return from json.Unmarshal is ignored
package main

import "encoding/json"

func main() {
    var data map[string]interface{}
    json.Unmarshal([]byte(`{"key":"value"}`), &data)
}
```

```golang
// Error return from http.ListenAndServe is ignored
package main

import "net/http"

func main() {
    http.ListenAndServe(":8080", nil)
}
```

### Valid

```golang
// Error from os.Open is properly checked
package main

import (
    "log"
    "os"
)

func main() {
    f, err := os.Open("file.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
}
```

```golang
// Error from json.Unmarshal is assigned and handled
package main

import (
    "encoding/json"
    "log"
)

func main() {
    var data map[string]interface{}
    err := json.Unmarshal([]byte(`{"key":"value"}`), &data)
    if err != nil {
        log.Fatal(err)
    }
}
```

```golang
// Error from http.ListenAndServe is properly handled
package main

import (
    "log"
    "net/http"
)

func main() {
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
```

```golang
// fmt.Println is excluded by default, so this is valid
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}
```

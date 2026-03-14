---
title: noBlankErrorAssignment
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noBlankErrorAssignment`
- This rule is **not recommended**, meaning it is not enabled by default. You can enable it manually.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noBlankErrorAssignment:
    enabled: true
```

## Details

Detects when error return values are explicitly assigned to the blank identifier (`_`). While assigning an error to `_` acknowledges that the function returns an error, it intentionally discards the error value, which can hide important failures.

This check is disabled by default because assigning to `_` is considered an explicit decision by the developer to ignore the error. However, enabling this rule helps enforce stricter error handling practices where all errors should be properly inspected, even if the developer is aware of the error return.

Common patterns flagged by this rule include:
- `_, _ = fmt.Fprintf(w, "...")` where the error is explicitly discarded
- `_ = f.Close()` where a file close error is discarded
- `_ = conn.SetDeadline(time.Now())` where a connection error is discarded

Source: https://github.com/kisielk/errcheck

## Examples

### Invalid

```golang
// Error from os.Open is explicitly discarded with blank identifier
package main

import "os"

func main() {
    f, _ := os.Open("file.txt")
    _ = f.Close()
}
```

```golang
// Error from json.NewEncoder Write is explicitly discarded
package main

import (
    "encoding/json"
    "os"
)

func main() {
    enc := json.NewEncoder(os.Stdout)
    _ = enc.Encode(map[string]string{"key": "value"})
}
```

```golang
// Error from database operation is discarded
package main

import "database/sql"

func closeDB(db *sql.DB) {
    _ = db.Close()
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
    err = f.Close()
    if err != nil {
        log.Printf("failed to close file: %v", err)
    }
}
```

```golang
// Error is assigned to a variable and handled
package main

import (
    "encoding/json"
    "log"
    "os"
)

func main() {
    enc := json.NewEncoder(os.Stdout)
    if err := enc.Encode(map[string]string{"key": "value"}); err != nil {
        log.Fatal(err)
    }
}
```

```golang
// Error from Close is captured and logged
package main

import (
    "database/sql"
    "log"
)

func closeDB(db *sql.DB) {
    if err := db.Close(); err != nil {
        log.Printf("failed to close database: %v", err)
    }
}
```

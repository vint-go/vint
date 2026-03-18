---
title: noDeferInInfiniteLoop
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noDeferInInfiniteLoop`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noDeferInInfiniteLoop:
    # rule options here
```

## Details

Defers in infinite loops will never execute.

Deferred function calls are executed when the surrounding function returns. In an infinite loop, the function never returns, so deferred calls will never execute and resources will leak.

Source: https://staticcheck.dev/docs/checks/#SA5003

## Examples

### Invalid

```golang
package main

import "os"

func process() {
    for {
        f, _ := os.Open("file.txt")
        // This defer will never execute
        defer f.Close()
        // process file...
    }
}
```

### Valid

```golang
package main

import "os"

func processFile() {
    f, _ := os.Open("file.txt")
    defer f.Close()
    // process file...
}

func process() {
    for {
        processFile()
    }
}
```

---
title: noEmptyForLoop
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noEmptyForLoop`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noEmptyForLoop:
    # rule options here
```

## Details

The empty `for` loop (`for {}`) spins and can consume 100% CPU. Use a select statement, channel, or `runtime.Gosched()` to yield.

An empty `for` loop without any blocking operations will busy-wait, consuming an entire CPU core. This is almost never intentional.

Source: https://staticcheck.dev/docs/checks/#SA5002

## Examples

### Invalid

```golang
package main

func main() {
    // Busy-wait loop consuming 100% CPU
    for {
    }
}
```

### Valid

```golang
package main

func main() {
    done := make(chan struct{})
    // Block efficiently on a channel
    select {
    case <-done:
    }
}
```

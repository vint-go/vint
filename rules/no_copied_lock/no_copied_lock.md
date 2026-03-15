---
title: noCopiedLock
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noCopiedLock`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noCopiedLock:
    # rule options here
```

## Details

Checks for locks erroneously passed by value. Values containing types from the `sync` package (such as `sync.Mutex`, `sync.RWMutex`, `sync.WaitGroup`, and `sync.Cond`) should not be copied, as copying may leave the lock in an inconsistent state. This analyzer detects:

- Function parameters that contain lock values (should be pointers)
- Assignments that copy lock values
- Return statements that copy lock values
- Range loop variables that copy lock values

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/copylock

## Examples

### Invalid

```golang
import "sync"

// Bad: mutex is copied when passed by value
func doSomething(mu sync.Mutex) {
    mu.Lock()
    defer mu.Unlock()
}
```

```golang
import "sync"

type MyStruct struct {
    mu sync.Mutex
}

func example() {
    a := MyStruct{}
    b := a // Bad: copies the mutex inside MyStruct
    _ = b
}
```

### Valid

```golang
import "sync"

// Good: mutex is passed by pointer
func doSomething(mu *sync.Mutex) {
    mu.Lock()
    defer mu.Unlock()
}
```

```golang
import "sync"

type MyStruct struct {
    mu sync.Mutex
}

func example() {
    a := &MyStruct{}
    b := a // Good: copies the pointer, not the mutex
    _ = b
}
```

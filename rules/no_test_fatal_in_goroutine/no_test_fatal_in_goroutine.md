---
title: noTestFatalInGoroutine
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noTestFatalInGoroutine`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noTestFatalInGoroutine:
    # rule options here
```

## Details

Detects calls to `Fatal`, `Fatalf`, `FailNow`, and similar methods from `testing.T` or `testing.B` from within a test goroutine other than the one that created the test. The `testing` package's fatal functions call `runtime.Goexit()` to stop the goroutine, but if called from a goroutine other than the test goroutine, it will not properly terminate the test and may cause unexpected behavior.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/testinggoroutine

## Examples

### Invalid

```golang
import "testing"

func TestExample(t *testing.T) {
    go func() {
        // Bad: t.Fatal called from a different goroutine
        t.Fatal("something failed")
    }()
}
```

```golang
import "testing"

func TestExample(t *testing.T) {
    done := make(chan bool)
    go func() {
        // Bad: t.Fatalf called from a spawned goroutine
        if err := doWork(); err != nil {
            t.Fatalf("work failed: %v", err)
        }
        done <- true
    }()
    <-done
}
```

### Valid

```golang
import "testing"

func TestExample(t *testing.T) {
    // Good: t.Fatal called from the test goroutine
    if err := doWork(); err != nil {
        t.Fatal(err)
    }
}
```

```golang
import "testing"

func TestExample(t *testing.T) {
    done := make(chan error)
    go func() {
        // Good: send the error back to the test goroutine
        done <- doWork()
    }()
    if err := <-done; err != nil {
        t.Fatal(err) // called from the test goroutine
    }
}
```

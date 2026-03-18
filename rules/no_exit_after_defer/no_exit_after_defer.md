---
title: noExitAfterDefer
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noExitAfterDefer`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noExitAfterDefer:
    # no additional options
```

## Details

Detects calls to `os.Exit`, `log.Fatal`, `log.Fatalf`, and `log.Fatalln` inside functions that use `defer`. These functions terminate the program immediately without running deferred cleanup operations, which is typically a bug. Deferred functions are only called when the surrounding function returns normally, so any resource cleanup registered with `defer` will be silently skipped.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
func main() {
    f, err := os.Open("file.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close() // will never execute if log.Fatal is called above
    // ...
    os.Exit(1) // deferred f.Close() will not execute
}
```

### Valid

```golang
func main() {
    if err := run(); err != nil {
        log.Fatal(err) // no defer in this function
    }
}

func run() error {
    f, err := os.Open("file.txt")
    if err != nil {
        return err
    }
    defer f.Close()
    // ...
    return nil
}
```

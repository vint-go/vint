---
title: noDeferInLoop
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noDeferInLoop`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noDeferInLoop:
    # no additional options
```

## Details

Detects defer statements inside loops. Deferred functions don't execute until the enclosing function returns, so using `defer` inside a loop causes resources to accumulate during loop iterations. This can exhaust memory or file handles. The fix is to refactor the loop body into a separate function where deferred cleanup executes after each iteration.

This rule is distinct from `lint/correctness/noDeferTimeMisuse`, which specifically detects misuse of `defer` with `time.Since` where the argument is evaluated immediately rather than at defer time.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
for _, filename := range filenames {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close() // deferred close won't execute until function returns
    // process file
}
```

### Valid

```golang
for _, filename := range filenames {
    if err := processFile(filename); err != nil {
        return err
    }
}

func processFile(filename string) error {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close() // properly deferred in a per-iteration function
    // process file
    return nil
}
```

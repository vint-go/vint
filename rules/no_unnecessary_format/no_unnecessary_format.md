---
title: noUnnecessaryFormat
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/noUnnecessaryFormat`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/noUnnecessaryFormat:
    # rule options here
```

## Details

This rule identifies calls to formatting functions where the format string does not contain any formatting verbs
and recommends switching to the non-formatting, more efficient alternative.

It checks functions from the `fmt`, `log`, `testing`, and `trace` packages, such as `fmt.Sprintf`, `fmt.Errorf`,
`log.Printf`, `t.Errorf`, and others. When the format string argument is a string literal that does not contain
a `%` character, the rule suggests using the simpler non-formatting variant (e.g., `fmt.Sprint` instead of
`fmt.Sprintf`, `errors.New` instead of `fmt.Errorf`).

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
fmt.Sprintf("hello world")      // use fmt.Sprint or just the string itself
fmt.Errorf("something failed")  // use errors.New
fmt.Printf("hello")             // use fmt.Print or fmt.Println
log.Printf("starting server")   // use log.Print
t.Errorf("test failed")         // use t.Error
```

### Valid

```golang
fmt.Sprintf("hello %s", name)
fmt.Errorf("failed: %w", err)
fmt.Printf("count: %d", n)
log.Printf("server on port %d", port)
t.Errorf("expected %v, got %v", want, got)
```

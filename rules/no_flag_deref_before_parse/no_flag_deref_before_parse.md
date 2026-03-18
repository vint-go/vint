---
title: noFlagDerefBeforeParse
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noFlagDerefBeforeParse`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noFlagDerefBeforeParse:
    # no additional options
```

## Details

Detects immediate dereferencing of `flag` package return values. Functions in the `flag` package (like `flag.String`, `flag.Int`, etc.) return pointers. Immediately dereferencing them before `flag.Parse()` is called results in getting the zero value, not the parsed value. The value should only be dereferenced after `flag.Parse()`.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
verbose := *flag.Bool("verbose", false, "enable verbose mode")
flag.Parse()
// verbose is always false here
```

### Valid

```golang
verbose := flag.Bool("verbose", false, "enable verbose mode")
flag.Parse()
if *verbose {
    // use after Parse
}
```

---
title: noOverwrittenArgument
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noOverwrittenArgument`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noOverwrittenArgument:
    # rule options here
```

## Details

A function argument is overwritten before first use.

When a function parameter is reassigned before it is ever read, the original argument value is discarded. This may indicate that the wrong variable name was used.

Source: https://staticcheck.dev/docs/checks/#SA4009

## Examples

### Invalid

```golang
package main

func process(name string) string {
    // name parameter is overwritten before use
    name = "default"
    return name
}
```

### Valid

```golang
package main

func process(name string) string {
    if name == "" {
        name = "default"
    }
    return name
}
```

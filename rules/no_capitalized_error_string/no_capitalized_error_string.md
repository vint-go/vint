---
title: noCapitalizedErrorString
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noCapitalizedErrorString`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noCapitalizedErrorString:
    # rule options here
```

## Details

Incorrectly formatted error string.

Error strings should not be capitalized (unless beginning with proper nouns or acronyms) and should not end with punctuation. This is because error strings are often composed with other messages.

Source: https://staticcheck.dev/docs/checks/#ST1005

## Examples

### Invalid

```golang
package main

import "errors"

func process() error {
    return errors.New("Something went wrong.")
}
```

### Valid

```golang
package main

import "errors"

func process() error {
    return errors.New("something went wrong")
}
```

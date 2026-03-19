---
title: noErrorStrings
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noErrorStrings`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noErrorStrings:
    arguments:
      - "xerrors.Errorf"
```

## Details

By convention, for better readability, error messages should not be capitalized or end with punctuation or a newline.

By default, the rule analyzes functions for creating errors from `fmt`, `errors`, and `github.com/pkg/errors`. Optionally, the rule can be configured to analyze user functions that create errors by passing additional function names in the `package.FunctionName` format.

More information: https://go.dev/wiki/CodeReviewComments#error-strings

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
package main

import "errors"

func example() error {
    // Error string starts with a capital letter
    return errors.New("Something went wrong")
}
```

```golang
package main

import "fmt"

func example() error {
    // Error string ends with punctuation
    return fmt.Errorf("something went wrong.")
}
```

```golang
package main

import "errors"

func example() error {
    // Error string ends with an exclamation mark
    return errors.New("something went wrong!")
}
```

### Valid

```golang
package main

import "errors"

func example() error {
    // Error string is lowercase and has no trailing punctuation
    return errors.New("something went wrong")
}
```

```golang
package main

import "fmt"

func example() error {
    // Error string starting with an acronym is acceptable
    return fmt.Errorf("HTTP request failed with status %d", code)
}
```

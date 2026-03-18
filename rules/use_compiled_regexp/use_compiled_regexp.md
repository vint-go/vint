---
title: useCompiledRegexp
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/useCompiledRegexp`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/useCompiledRegexp:
    # rule options here
```

## Details

Using `regexp.Match` or related functions in a loop, should use `regexp.Compile`.

Calling `regexp.MatchString` or similar functions inside a loop recompiles the regular expression on every iteration. Compile the regex once outside the loop and reuse it for better performance.

Source: https://staticcheck.dev/docs/checks/#SA6000

## Examples

### Invalid

```golang
package main

import "regexp"

func containsDigits(items []string) []bool {
    results := make([]bool, len(items))
    for i, item := range items {
        // Regex is recompiled on each iteration
        results[i], _ = regexp.MatchString(`\d+`, item)
    }
    return results
}
```

### Valid

```golang
package main

import "regexp"

var digitRe = regexp.MustCompile(`\d+`)

func containsDigits(items []string) []bool {
    results := make([]bool, len(items))
    for i, item := range items {
        results[i] = digitRe.MatchString(item)
    }
    return results
}
```

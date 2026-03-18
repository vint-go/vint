---
title: noInvalidTemplate
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInvalidTemplate`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInvalidTemplate:
    # rule options here
```

## Details

Invalid template.

This check validates that arguments to `template.Must`, `template.New().Parse()`, and related functions contain valid Go templates. Invalid templates will cause runtime errors.

Source: https://staticcheck.dev/docs/checks/#SA1001

## Examples

### Invalid

```golang
package main

import "text/template"

func main() {
    // Invalid template - unclosed action
    t := template.Must(template.New("test").Parse("{{.Name"))
    _ = t
}
```

### Valid

```golang
package main

import "text/template"

func main() {
    // Valid template
    t := template.Must(template.New("test").Parse("{{.Name}}"))
    _ = t
}
```

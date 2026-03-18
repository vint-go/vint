---
title: noSelfReferencingFinalizer
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSelfReferencingFinalizer`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSelfReferencingFinalizer:
    # rule options here
```

## Details

The finalizer references the finalized object, preventing garbage collection.

When using `runtime.SetFinalizer`, the finalizer function must not reference the object being finalized through a closure, as this creates a reference cycle that prevents the object from ever being garbage collected.

Source: https://staticcheck.dev/docs/checks/#SA5005

## Examples

### Invalid

```golang
package main

import "runtime"

type Resource struct {
    name string
}

func newResource(name string) *Resource {
    r := &Resource{name: name}
    // Closure captures r, preventing GC
    runtime.SetFinalizer(r, func(_ *Resource) {
        cleanup(r)
    })
    return r
}

func cleanup(r *Resource) {}
```

### Valid

```golang
package main

import "runtime"

type Resource struct {
    name string
}

func newResource(name string) *Resource {
    r := &Resource{name: name}
    // Use the parameter, not the closure variable
    runtime.SetFinalizer(r, func(r *Resource) {
        cleanup(r)
    })
    return r
}

func cleanup(r *Resource) {}
```

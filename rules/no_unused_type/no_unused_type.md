---
title: noUnusedType
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnusedType`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnusedType:
    # Controls whether types in generated files are considered used.
    # generated-is-used: true
```

## Details

Detects named types that are declared but never used anywhere in the codebase.

The `unused` linter builds a dependency graph to track how types are referenced. A named type is considered "used" if it is reachable from an entry point. This includes usage in variable declarations, function signatures, type assertions, type conversions, composite literals, and embedded fields.

Key rules governing type usage:
- Exported named types are always considered used (since external packages may reference them).
- Named types use their exported methods (rule 2.1), meaning if a type is used, its exported methods are also considered used.
- Named types use the type they are based on (rule 2.2), as well as all their type parameters (rule 2.5) and type arguments (rule 2.6).
- Types use their underlying and element types (rule 9.3).
- Conversions use the type they convert to (rule 9.4).
- Type parameters use their constraint type (rule 12.1).
- Types in generated files are considered used when the `generated-is-used` option is enabled (default).

Note: Type parameters (generics constraints like `T` in `func Foo[T any]()`) are detected by the underlying analyzer but are filtered out and not reported by the golangci-lint integration.

Source: https://github.com/dominikh/go-tools/tree/master/unused

## Examples

### Invalid

```golang
package mypackage

type config struct { // type config is unused
    host string
    port int
}

func Connect() {}
```

```golang
package mypackage

type Result struct { // type Result is unused (unexported in non-main package context)
    Value int
    Error error
}

// Note: if this were "type result struct", and it's never used anywhere, it would be reported.
```

```golang
package mypackage

type handler func(req string) string // type handler is unused

func Process(input string) string {
    return input
}
```

```golang
package mypackage

type Color int // type Color is unused

const (
    Red   Color = iota
    Green
    Blue
)
```

### Valid

```golang
package mypackage

type config struct {
    host string
    port int
}

func NewConfig() config {
    return config{host: "localhost", port: 8080}
}
```

```golang
package mypackage

// Exported types are always considered used
type Server struct {
    Addr string
}
```

```golang
package mypackage

type stringer interface {
    String() string
}

func Format(s stringer) string {
    return s.String()
}
```

```golang
package mypackage

type middleware func(next handler) handler

type handler func(req string) string

func Chain(h handler, mws ...middleware) handler {
    for _, mw := range mws {
        h = mw(h)
    }
    return h
}
```

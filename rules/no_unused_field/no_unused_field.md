---
title: noUnusedField
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnusedField`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnusedField:
    # Controls whether writing to a field counts as using it.
    # field-writes-are-uses: true
    # Controls whether exported struct fields are always considered used.
    # exported-fields-are-used: true
    # Controls whether fields in generated files are considered used.
    # generated-is-used: true
```

## Details

Detects struct fields that are declared but never accessed (read or written to) anywhere in the codebase.

The `unused` linter uses a graph-based approach to determine whether struct fields are reachable. A field is considered "used" if it is accessed (via read or write) in code that is itself reachable from an entry point.

Key rules governing field usage:
- Exported fields are considered used by default when the `exported-fields-are-used` option is enabled (default). When disabled, exported fields are also subject to unused detection.
- Writing to a field counts as using it when `field-writes-are-uses` is enabled (default). Disabling this option means that only reading from a field counts as using it.
- Fields of type `NoCopy` sentinel are always considered used (rule 6.1).
- Embedded fields that help implement interfaces are considered used (rule 6.3, recursive).
- Embedded fields that have exported methods are considered used (rule 6.4, recursive).
- Embedded structs that have exported fields are considered used (rule 6.5, recursive).
- All fields are considered used if the struct has a `structs.HostLayout` field (rule 6.6).
- When converting between two equivalent structs, fields use each other (rule 5.1).
- When converting to or from `unsafe.Pointer`, all fields are marked as used (rule 5.2).
- Fields of anonymous struct types are always considered used (rule 11.1).
- Field accesses use their corresponding fields (rule 7.1).
- Fields use their types (rule 7.2).

Source: https://github.com/dominikh/go-tools/tree/master/unused

## Examples

### Invalid

```golang
package mypackage

type User struct {
    Name     string
    Age      int
    nickname string // field nickname is unused
}

func NewUser(name string, age int) User {
    return User{Name: name, Age: age}
}

func (u User) Display() string {
    return u.Name
}
```

```golang
package mypackage

type config struct {
    host    string
    port    int
    timeout int // field timeout is unused
}

func newConfig() config {
    return config{host: "localhost", port: 8080}
}

func (c config) Address() string {
    return c.host + ":" + string(rune(c.port))
}
```

```golang
package mypackage

type response struct {
    status  int
    body    string
    headers map[string]string // field headers is unused
}

func ok(body string) response {
    return response{status: 200, body: body}
}

func (r response) String() string {
    return r.body
}
```

### Valid

```golang
package mypackage

type User struct {
    Name     string
    Age      int
    nickname string
}

func NewUser(name string, age int, nick string) User {
    return User{Name: name, Age: age, nickname: nick}
}

func (u User) Nickname() string {
    return u.nickname // nickname is read here
}
```

```golang
package mypackage

// Exported fields are considered used by default
type Config struct {
    Host    string
    Port    int
    Timeout int
}
```

```golang
package mypackage

import "sync"

type SafeMap struct {
    mu sync.Mutex // embedded/NoCopy-style fields are considered used
    m  map[string]string
}

func (s *SafeMap) Set(k, v string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.m[k] = v
}

func (s *SafeMap) Get(k string) string {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.m[k]
}
```

```golang
package mypackage

type Reader interface {
    Read(p []byte) (n int, err error)
}

type MyReader struct {
    // This embedded field helps implement the Reader interface,
    // so it is considered used (rule 6.3)
    innerReader Reader
}
```

```golang
package mypackage

// Anonymous struct fields are always considered used (rule 11.1)
func process() {
    data := struct {
        Name  string
        Value int
    }{
        Name:  "test",
        Value: 42,
    }
    _ = data.Name
}
```

---
title: noImpossibleInterfaceAssert
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noImpossibleInterfaceAssert`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noImpossibleInterfaceAssert:
    # rule options here
```

## Details

Flags impossible interface-to-interface type assertions. This analyzer detects type assertions from one interface type to another where the assertion can never succeed because the two interfaces contain methods with the same name but different signatures. Since no type could implement both interfaces simultaneously, the assertion is guaranteed to fail.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/ifaceassert

## Examples

### Invalid

```golang
type Reader interface {
    Read(p []byte) (int, error)
}

type BadReader interface {
    Read(p []byte) error // different signature than Reader.Read
}

func example(r Reader) {
    // Impossible: no type can have Read with both signatures
    _ = r.(BadReader)
}
```

### Valid

```golang
type Reader interface {
    Read(p []byte) (int, error)
}

type ReadCloser interface {
    Read(p []byte) (int, error)
    Close() error
}

func example(r Reader) {
    // Possible: ReadCloser extends Reader
    if rc, ok := r.(ReadCloser); ok {
        defer rc.Close()
    }
}
```

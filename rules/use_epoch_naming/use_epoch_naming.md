---
title: useEpochNaming
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useEpochNaming`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useEpochNaming:
    # no options
```

## Details

Variables initialized with epoch time methods (`time.Now().Unix()`, `time.Now().UnixMilli()`,
`time.Now().UnixMicro()`, `time.Now().UnixNano()`) should have names that clearly indicate their time unit to
prevent confusion and potential bugs when working with different time scales.

This rule enforces that variable names contain appropriate suffixes based on the method used:

- `Unix()`: variable name should end with "Sec", "Second" or "Seconds"
- `UnixMilli()`: variable name should end with "Milli" or "Ms"
- `UnixMicro()`: variable name should end with "Micro", "Microsecond", "Microseconds" or "Us"
- `UnixNano()`: variable name should end with "Nano" or "Ns"

The rule checks variable declarations, short variable declarations (`:=`), and regular assignments (`=`).
The suffix matching is case-insensitive and must appear at the end of the variable name.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
timestamp := time.Now().Unix()           // unclear which unit
createdAt := time.Now().UnixMilli()      // missing unit indicator
t := time.Now().UnixNano()               // lacks required suffix
```

### Valid

```golang
timestampSec := time.Now().Unix()        // clearly seconds
createdAtMs := time.Now().UnixMilli()    // clearly milliseconds
tNano := time.Now().UnixNano()           // clearly nanoseconds
createdSeconds := time.Now().Unix()      // full word is fine
updatedMicro := time.Now().UnixMicro()   // microseconds
```

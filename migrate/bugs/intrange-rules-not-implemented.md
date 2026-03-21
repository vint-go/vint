# intrange: vint rules not yet implemented

## Affected rules
- `lint/style/noRedundantRangeLen` (maps to golangci-lint `intrange` / `simplify-range-len`)
- `lint/style/useIntegerRange` (maps to golangci-lint `intrange` / `use-integer-range`)

## Behavior in golangci-lint
The intrange linter detects two patterns:
1. `for i := range len(slice)` which can be simplified to `for i := range slice`
2. C-style `for i := 0; i < n; i++` loops that can use Go 1.22+ integer range syntax (`for i := range n`)

## Behavior in vint
The rules `noRedundantRangeLen` and `useIntegerRange` are registered in the rules registry but do not have Go rule implementations in the `rules/` directory. The migrator correctly maps configuration, but the rules cannot fire during nolint conversion.

## Gap
- Config migration works correctly: enabling `intrange` in golangci-lint produces the correct `vint.yaml` with both rules enabled.
- Nolint conversion is degraded: `//nolint:intrange` directives are replaced with a "no vint rules fired" note instead of proper `// vint-ignore` directives, because the rules have no implementation to detect violations.

## Example
Input:
```go
for i := 0; i < 10; i++ { //nolint:intrange
    fmt.Println(i)
}
```

Expected (once rules are implemented):
```go
// vint-ignore lint/style/useIntegerRange: migrated from nolint:intrange
for i := 0; i < 10; i++ {
    fmt.Println(i)
}
```

Actual:
```go
// NOTE: nolint removed — no vint rules fired (was: //nolint:intrange)
for i := 0; i < 10; i++ {
    fmt.Println(i)
}
```

## Impact on migration
Users migrating from golangci-lint will have their `//nolint:intrange` directives removed with a note rather than converted to `// vint-ignore` directives. Once the vint rules are implemented, the nolint test should be updated to verify proper conversion.

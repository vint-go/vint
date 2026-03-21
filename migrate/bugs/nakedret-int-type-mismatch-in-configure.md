# nakedret: max-func-lines setting ignored during nolint conversion due to int type mismatch

## Affected rule
`lint/style/noNakedReturn` (maps to golangci-lint `nakedret`)

## Behavior in golangci-lint
When `max-func-lines` is set (e.g., to 10), nakedret only flags naked returns in functions longer than 10 lines.

## Behavior in vint
The `noNakedReturn` rule's `Configure` method type-asserts the `maxFuncLines` value as `int64`:
```go
lines, ok := v.(int64)
```
However, when settings flow through the migration system, YAML v3 decodes integers as Go `int`, not `int64`. The type assertion silently fails and the rule falls back to its default of 30 lines.

## Gap
When the migrator passes `max-func-lines` through to the vint rule during nolint conversion, the configured threshold is silently ignored. The rule runs with the default threshold of 30 lines instead. This means:
- Functions between the user's configured threshold and 30 lines would not have their `//nolint:nakedret` directives properly converted to `// vint-ignore` comments.
- The generated `vint.yaml` correctly contains the setting, but the nolint conversion step uses the wrong threshold.

## Example
With golangci-lint config `max-func-lines: 10`, a 15-line function with a naked return and `//nolint:nakedret` would not get its nolint converted to vint-ignore because the rule runs with threshold 30 (not 10) during conversion.

## Impact on migration
Low impact for most users since the default threshold of 30 is the same in both tools. Users who have customized `max-func-lines` to a lower value may see some `//nolint:nakedret` directives converted to "no vint rules fired" notes instead of proper `// vint-ignore` comments. The generated `vint.yaml` config file itself is correct.

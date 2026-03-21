# govet: per-analyzer settings not migrated

## Affected rules
- `lint/correctness/noPrintfFormatMismatch` (maps to golangci-lint `govet` analyzer `printf`)
- `lint/suspicious/noVariableShadowing` (maps to golangci-lint `govet` analyzer `shadow`)
- `lint/suspicious/noSpecificFunctionCall` (maps to golangci-lint `govet` analyzer `findcall`)

## Behavior in golangci-lint
golangci-lint's `govet` linter supports per-analyzer settings via the `settings` map. For example:
- `printf.funcs`: additional function names to check for printf-style format strings
- `shadow.strict`: enables strict mode for variable shadowing detection

## Behavior in vint
The vint rules `noPrintfFormatMismatch` and `noVariableShadowing` do not have `Configure()` methods, so they cannot accept per-analyzer settings. The migrator passes through `findcall.name` since `noSpecificFunctionCall` does support configuration.

## Gap
Per-analyzer settings for `printf.funcs` and `shadow.strict` are silently dropped during migration. Users who relied on custom printf function lists or strict shadow checking will get different behavior after migration.

## Example
```yaml
# golangci-lint config
linters:
  settings:
    govet:
      settings:
        printf:
          funcs:
            - (github.com/my/pkg.Logger).Infof
        shadow:
          strict: true
```

These settings have no equivalent in vint and are lost during migration.

## Impact on migration
Users with custom `printf.funcs` will not get format string checking on their custom logging functions. Users with `shadow.strict: true` may miss some variable shadowing warnings they previously received.

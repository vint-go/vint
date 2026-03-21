# errorlint: allowed-errors settings not migrated

## Affected rule
`lint/correctness/noDirectErrorComparison` (maps to golangci-lint `errorlint` comparison check)

## Behavior in golangci-lint
errorlint supports `allowed-errors` and `allowed-errors-wildcard` settings that let users whitelist specific error/function pairs where direct error comparison is acceptable:

```yaml
linters:
  settings:
    errorlint:
      allowed-errors:
        - err: "io.EOF"
          fun: "example.com/pkg.Read"
      allowed-errors-wildcard:
        - err: "example.com/pkg.ErrMagic"
          fun: "example.com/pkg.*"
```

## Behavior in vint
The `noDirectErrorComparison` rule does not support `allowed-errors` or `allowed-errors-wildcard` configuration options. It has a hardcoded exception for `io.EOF` but no mechanism for user-defined allowed error/function pairs.

## Gap
Users who have configured `allowed-errors` or `allowed-errors-wildcard` in their golangci-lint config will lose those exceptions after migration. The migrator currently does not pass these settings through because the vint rule has no corresponding configuration.

## Example
```yaml
# golangci-lint config that cannot be fully migrated
linters:
  settings:
    errorlint:
      allowed-errors:
        - err: "example.com/pkg.ErrExpected"
          fun: "example.com/pkg.Process"
```

## Impact on migration
- Users with `allowed-errors` settings will get new warnings after migration on code that was previously allowed.
- The migration is safe (no suppressions are lost), but coverage is stricter than the original configuration.

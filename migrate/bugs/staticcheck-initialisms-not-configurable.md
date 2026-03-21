# staticcheck: custom initialisms not configurable

## Affected rule
`lint/style/useIdiomaticNaming` (maps to golangci-lint staticcheck check `ST1003`)

## Behavior in golangci-lint
The staticcheck linter in golangci-lint supports a `initialisms` setting that allows users to customize the list of initialisms (e.g. `ID`, `URL`, `HTTP`) used when checking naming conventions. Users can add domain-specific initialisms or remove defaults.

## Behavior in vint
The vint `useIdiomaticNaming` rule does not have a `Configure` method and does not accept custom initialisms. It uses a hardcoded set of initialisms via `internalrule.Name()`.

## Gap
Custom initialisms from golangci-lint staticcheck configuration cannot be migrated to vint. The migration drops the `initialisms` setting silently.

## Example
```yaml
# golangci-lint config
linters:
  settings:
    staticcheck:
      initialisms:
        - ID
        - URL
        - HTTP
        - GRPC  # custom
```

With the custom `GRPC` initialism, staticcheck would flag `GrpcServer` and suggest `GRPCServer`. Vint's rule would not recognize `GRPC` as a custom initialism and would not produce the same suggestion.

## Impact on migration
Users who have customized the `initialisms` list will get the default initialism behavior after migration. This may produce different naming suggestions than expected. The `http-status-code-whitelist` setting is also not passed through, though its impact is limited to the `ST1013` check which is not currently implemented in vint.

# gosec: Per-rule configuration not migrated

## Affected rules
All gosec rules in `lint/security/*`

## Behavior in golangci-lint
gosec supports detailed per-rule configuration via the `config` map in linter settings:
- `global.nosec`: Controls whether `#nosec` annotations are honored
- `global.audit`: Enables audit mode with stricter checks
- `G101.pattern`: Custom regex pattern for credential name matching
- `G101.entropy_threshold`: Entropy threshold for string analysis
- `G301`, `G302`, `G306`: Custom permission thresholds (e.g., `"0750"`)
- `G104.fmt`: List of functions to check for unchecked errors
- `G111.pattern`: Custom pattern for directory serving detection

Additionally, `severity` and `confidence` settings filter which findings are reported.

## Behavior in vint
The vint rule implementations for gosec rules do not accept configuration arguments. All rules use hardcoded defaults. The migrator does not pass through any per-rule configuration.

## Gap
Users who have customized gosec behavior through the `config` map will not have those customizations reflected in the migrated vint configuration. The `severity` and `confidence` filters are also not applicable in vint.

## Example
```yaml
# golangci-lint config
linters:
  settings:
    gosec:
      config:
        G101:
          pattern: "(?i)my_custom_secret"
          entropy_threshold: "80.0"
        G301: "0750"
      severity: high
      confidence: medium
```
These settings would be silently dropped during migration.

## Impact on migration
Medium impact for users with customized gosec settings. Projects using default gosec configuration (the common case) are unaffected. Users with custom patterns or thresholds should manually verify that vint's default behavior covers their security requirements.

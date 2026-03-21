# gocritic: enabled/disabled checks configuration not migrated

## Affected rule
All gocritic-mapped vint rules

## Behavior in golangci-lint
gocritic supports `enabled-checks`, `disabled-checks`, `enabled-tags`, and `disabled-tags` configuration options that allow users to selectively enable or disable specific checks or groups of checks. For example, a user might disable the `commentedOutCode` check or only enable checks tagged as `diagnostic`.

## Behavior in vint
The migrator unconditionally enables all 95 mapped vint rules when gocritic is enabled, regardless of any `enabled-checks`, `disabled-checks`, `enabled-tags`, or `disabled-tags` settings in the golangci-lint configuration.

## Gap
Users who had selectively enabled or disabled specific gocritic checks will get all checks enabled after migration. This may produce more (or fewer) warnings than expected.

## Example
```yaml
# golangci-lint config
linters:
  enable:
    - gocritic
  settings:
    gocritic:
      disabled-checks:
        - commentedOutCode
        - hugeParam
```

After migration, both `lint/style/noCommentedOutCode` and `lint/performance/noHugeParam` will be enabled even though the user had disabled them.

## Impact on migration
Users with customized gocritic check lists will need to manually disable unwanted rules in vint.yaml after migration. The migrator produces a superset of what the user intended.

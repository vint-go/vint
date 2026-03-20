# golangci-lint to Vint migrator

This tool is designed to assist in migrating from golangci-lint to Vint. It provides a structured approach to transition your linting configuration and rules, ensuring a smooth and efficient migration process.

## Features

- **Configuration Mapping**: Automatically translates golangci-lint configuration to Vint's format.
- **Rule Coverage Check**: Identifies if all golangci-lint rules used in your configuration are supported by Vint.
- **nolint Directive Handling**: Converts golangci-lint's `//nolint` directives to Vint's equivalent, ensuring that your code annotations remain effective.

## How it works

### Configuration Mapping

The migrator separately understands every confuguration for linters.

File structure:

```
migrator/
├── linters/
│   ├── gocritic/
│   ├── staticcheck/
│   └── ... # all linters supported by migrator
```

### Rule Coverage Check

Migrator has access to mapping from linter to known rule. If rule is not implemented in vint, migrator knows this and can report it to user. Linter might also be unknown to migrator, in which case migrator will report that it doesn't know this linter and it means it cannot be migrated.

### nolint Directive Handling

'nolint' directive disables the entire linter, but it is not obvious which rule it is there for. To cover for this, migrator will run all rules corresponding to linter on the file with 'nolint' directive and then inject vint's own suppression for specific rule (or rules) instead of for the whole linter.


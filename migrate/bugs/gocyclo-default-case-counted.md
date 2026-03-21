# gocyclo: default case counted in complexity

## Affected rule
`lint/complexity/noHighCyclomaticComplexity` (maps to golangci-lint `gocyclo`)

## Behavior in golangci-lint
The original gocyclo linter (github.com/fzipp/gocyclo) excludes `default` cases from the complexity count. Specifically:
- `*ast.CaseClause` is only counted when `n.List != nil` (has explicit case values)
- `*ast.CommClause` is only counted when `n.Comm != nil` (has explicit comm statement)

This means `default:` branches in switch/select statements do not add to cyclomatic complexity.

## Behavior in vint
The vint rule counts ALL `*ast.CaseClause` and `*ast.CommClause` nodes, including `default:` cases, incrementing complexity for each one.

## Gap
Vint reports a slightly higher cyclomatic complexity than gocyclo for functions that contain `default:` branches in switch or select statements. Each `default:` case adds +1 to the complexity in vint but not in the original gocyclo.

## Example
```go
func example(x int) string {
    switch x {
    case 1:
        return "one"
    case 2:
        return "two"
    default:
        return "other"
    }
}
```

- gocyclo complexity: 3 (1 base + 2 case clauses, default excluded)
- vint complexity: 4 (1 base + 2 case clauses + 1 default)

## Impact on migration
Users migrating from golangci-lint may see additional complexity warnings on functions that were previously under their configured threshold. Functions with switch/select statements containing `default:` branches will report higher complexity in vint. This could cause previously-passing functions to fail the complexity check, especially when the function's complexity is near the threshold.

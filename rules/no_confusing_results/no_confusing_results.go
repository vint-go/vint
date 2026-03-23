package no_confusing_results

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// ConfusingResultsRule lints given function declarations.
type ConfusingResultsRule struct{}

// Apply applies the rule to given file.
func (*ConfusingResultsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)

		isFunctionWithMoreThanOneResult := ok && funcDecl.Type.Results != nil && len(funcDecl.Type.Results.List) > 1
		if !isFunctionWithMoreThanOneResult {
			continue
		}

		resultsAreNamed := len(funcDecl.Type.Results.List[0].Names) > 0
		if resultsAreNamed {
			continue
		}

		lastType := ""
		for _, result := range funcDecl.Type.Results.List {
			resultTypeName := astutils.GoFmt(result.Type)

			if resultTypeName == lastType {
				failures = append(failures, lint.Failure{
					Node:       result,
					Confidence: 1,
					Category:   lint.FailureCategoryNaming,
					Failure:    "unnamed results of the same type may be confusing, consider using named results",
				})

				break
			}

			lastType = resultTypeName
		}
	}

	return failures
}

// Name returns the rule name.
func (*ConfusingResultsRule) Name() string {
	return "noConfusingResults"
}

// Group returns the rule group.
func (*ConfusingResultsRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*ConfusingResultsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

package use_errors_new

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseErrorsNewRule spots calls to [fmt.Errorf] that can be replaced by [errors.New].
type UseErrorsNewRule struct{}

// Apply applies the rule to given file.
func (*UseErrorsNewRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	walker := lintFmtErrorf{
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, file.AST)

	return failures
}

// Name returns the rule name.
func (*UseErrorsNewRule) Name() string {
	return "useErrorsNew"
}

// Group returns the rule group.
func (*UseErrorsNewRule) Group() string {
	return "style"
}

type lintFmtErrorf struct {
	onFailure func(lint.Failure)
}

func (w lintFmtErrorf) Visit(n ast.Node) ast.Visitor {
	funcCall, ok := n.(*ast.CallExpr)
	if !ok {
		return w // not a function call
	}

	isFmtErrorf := astutils.IsPkgDotName(funcCall.Fun, "fmt", "Errorf")
	if !isFmtErrorf {
		return w // not a call to fmt.Errorf
	}

	if len(funcCall.Args) > 1 {
		return w // the use of fmt.Errorf is legit
	}

	// the call is of the form fmt.Errorf("...")
	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryErrors,
		Node:       n,
		Confidence: 1,
		Failure:    "replace fmt.Errorf by errors.New",
	})

	return w
}

// CacheTier returns the cache tier for this rule.
func (*UseErrorsNewRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

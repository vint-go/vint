package no_forbidden_call_in_wg_go

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// ForbiddenCallInWgGoRule spots calls to panic or wg.Done when using [sync.WaitGroup.Go].
type ForbiddenCallInWgGoRule struct{}

// Apply applies the rule to given file.
func (*ForbiddenCallInWgGoRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !file.Pkg.IsAtLeastGoVersion(lint.Go125) {
		return nil // skip analysis if Go version < 1.25
	}

	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintForbiddenCallInWgGo{
		onFailure: onFailure,
	}

	// Iterate over declarations looking for function declarations
	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue // not a function
		}

		if fn.Body == nil {
			continue // external (no-Go) function
		}

		// Analyze the function body
		ast.Walk(w, fn.Body)
	}

	return failures
}

// Name returns the rule name.
func (*ForbiddenCallInWgGoRule) Name() string {
	return "noForbiddenCallInWgGo"
}

// Group returns the rule group.
func (*ForbiddenCallInWgGoRule) Group() string {
	return "correctness"
}

type lintForbiddenCallInWgGo struct {
	onFailure func(lint.Failure)
}

func (w *lintForbiddenCallInWgGo) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w // not a call of statements
	}

	if !astutils.IsPkgDotName(call.Fun, "wg", "Go") {
		return w // not a call to wg.Go
	}

	if len(call.Args) != 1 {
		return nil // no argument (impossible)
	}

	funcLit, ok := call.Args[0].(*ast.FuncLit)
	if !ok {
		return nil // the argument is not a function literal
	}

	var callee string

	forbiddenCallPicker := func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return false
		}

		if astutils.IsPkgDotName(call.Fun, "wg", "Done") ||
			astutils.IsIdent(call.Fun, "panic") ||
			astutils.IsPkgDotName(call.Fun, "log", "Panic") ||
			astutils.IsPkgDotName(call.Fun, "log", "Panicf") ||
			astutils.IsPkgDotName(call.Fun, "log", "Panicln") {
			callee = astutils.GoFmt(n)
			callee, _, _ = strings.Cut(callee, "(")
			return true
		}

		return false
	}

	// search a forbidden call in the body of the function literal
	forbiddenCall := astutils.SeekNode[*ast.CallExpr](funcLit.Body, forbiddenCallPicker)
	if forbiddenCall == nil {
		return nil // there is no forbidden call in the call to wg.Go
	}

	msg := fmt.Sprintf("do not call %s inside wg.Go", callee)
	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       forbiddenCall,
		Category:   lint.FailureCategoryErrors,
		Failure:    msg,
	})

	return nil
}

// CacheTier returns the cache tier for this rule.
func (*ForbiddenCallInWgGoRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

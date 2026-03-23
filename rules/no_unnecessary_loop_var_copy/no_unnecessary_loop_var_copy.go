package no_unnecessary_loop_var_copy

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnnecessaryLoopVarCopyRule detects unnecessary copies of loop variables inside loop bodies.
// Starting with Go 1.22, each iteration of a for loop gets its own copy of the loop variable,
// making the common `v := v` idiom unnecessary.
type NoUnnecessaryLoopVarCopyRule struct {
	checkAlias bool
}

// Configure validates the rule configuration, and configures the rule accordingly.
func (r *NoUnnecessaryLoopVarCopyRule) Configure(arguments lint.Arguments) error {
	r.checkAlias = false
	if len(arguments) == 0 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noUnnecessaryLoopVarCopy" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		if k == "check-alias" || k == "checkAlias" || k == "checkalias" {
			val, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for check-alias in "noUnnecessaryLoopVarCopy" rule; need bool but got %T`, v)
			}
			r.checkAlias = val
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoUnnecessaryLoopVarCopyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoUnnecessaryLoopVarCopy{
		checkAlias: r.checkAlias,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules
func (r *NoUnnecessaryLoopVarCopyRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoUnnecessaryLoopVarCopy{
		checkAlias: r.checkAlias,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUnnecessaryLoopVarCopyRule) Name() string {
	return "noUnnecessaryLoopVarCopy"
}

// Group returns the rule group.
func (*NoUnnecessaryLoopVarCopyRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnnecessaryLoopVarCopyRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoUnnecessaryLoopVarCopy struct {
	checkAlias bool
	onFailure  func(lint.Failure)
}

func (w *lintNoUnnecessaryLoopVarCopy) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.RangeStmt:
		w.checkRangeStmt(n)
		return nil // we manually walk the body, avoid double-visiting
	case *ast.ForStmt:
		w.checkForStmt(n)
		return nil // we manually walk the body, avoid double-visiting
	}
	return w
}

// checkRangeStmt checks a range loop for unnecessary copies of the key/value variables.
func (w *lintNoUnnecessaryLoopVarCopy) checkRangeStmt(rangeStmt *ast.RangeStmt) {
	if rangeStmt.Tok != token.DEFINE {
		return
	}

	// Collect loop variable names
	loopVars := make(map[string]bool)
	if ident, ok := rangeStmt.Key.(*ast.Ident); ok && ident.Name != "_" {
		loopVars[ident.Name] = true
	}
	if rangeStmt.Value != nil {
		if ident, ok := rangeStmt.Value.(*ast.Ident); ok && ident.Name != "_" {
			loopVars[ident.Name] = true
		}
	}

	if len(loopVars) == 0 || rangeStmt.Body == nil {
		return
	}

	w.checkBody(rangeStmt.Body, loopVars)
}

// checkForStmt checks a C-style for loop for unnecessary copies of the init variables.
func (w *lintNoUnnecessaryLoopVarCopy) checkForStmt(forStmt *ast.ForStmt) {
	if forStmt.Init == nil || forStmt.Body == nil {
		return
	}

	// Extract variable names from the init statement
	loopVars := make(map[string]bool)
	assignStmt, ok := forStmt.Init.(*ast.AssignStmt)
	if !ok || assignStmt.Tok != token.DEFINE {
		return
	}

	for _, lhs := range assignStmt.Lhs {
		if ident, ok := lhs.(*ast.Ident); ok && ident.Name != "_" {
			loopVars[ident.Name] = true
		}
	}

	if len(loopVars) == 0 {
		return
	}

	w.checkBody(forStmt.Body, loopVars)
}

// checkBody scans the body of a for loop for short variable declarations that copy loop variables.
func (w *lintNoUnnecessaryLoopVarCopy) checkBody(body *ast.BlockStmt, loopVars map[string]bool) {
	// Create a new walker for the body that will check for loop var copies
	// but also continue walking nested structures
	bodyWalker := &loopBodyWalker{
		loopVars:     loopVars,
		checkAlias:   w.checkAlias,
		onFailure:    w.onFailure,
		parentWalker: w,
	}
	for _, stmt := range body.List {
		ast.Walk(bodyWalker, stmt)
	}
}

type loopBodyWalker struct {
	loopVars     map[string]bool
	checkAlias   bool
	onFailure    func(lint.Failure)
	parentWalker *lintNoUnnecessaryLoopVarCopy
}

func (bw *loopBodyWalker) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	// If we encounter a nested for loop, let the parent walker handle it
	// (it will extract its own loop variables)
	switch n := node.(type) {
	case *ast.RangeStmt:
		bw.parentWalker.checkRangeStmt(n)
		return nil
	case *ast.ForStmt:
		bw.parentWalker.checkForStmt(n)
		return nil
	case *ast.AssignStmt:
		if n.Tok == token.DEFINE {
			bw.checkAssignment(n)
		}
	}

	return bw
}

// checkAssignment checks a short variable declaration for unnecessary copies of loop variables.
func (bw *loopBodyWalker) checkAssignment(assign *ast.AssignStmt) {
	// Look at each RHS expression to see if it's a simple identifier that matches a loop variable
	for i, rhs := range assign.Rhs {
		rhsIdent, ok := rhs.(*ast.Ident)
		if !ok {
			continue
		}
		if !bw.loopVars[rhsIdent.Name] {
			continue
		}

		// The RHS is a loop variable. Now check the corresponding LHS.
		if i >= len(assign.Lhs) {
			continue
		}
		lhsIdent, ok := assign.Lhs[i].(*ast.Ident)
		if !ok {
			continue
		}

		// Check if this is a same-name copy (e.g., `i := i`)
		if lhsIdent.Name == rhsIdent.Name {
			bw.onFailure(lint.Failure{
				Confidence: 1,
				Category:   lint.FailureCategoryStyle,
				Node:       assign,
				Failure:    fmt.Sprintf(`The copy of the 'for' variable "%s" can be deleted (Go 1.22+)`, rhsIdent.Name),
			})
			continue
		}

		// Check if this is an alias copy (e.g., `_i := i`) and check-alias is enabled
		if bw.checkAlias {
			bw.onFailure(lint.Failure{
				Confidence: 1,
				Category:   lint.FailureCategoryStyle,
				Node:       assign,
				Failure:    fmt.Sprintf(`The copy of the 'for' variable "%s" can be deleted (Go 1.22+)`, rhsIdent.Name),
			})
		}
	}
}

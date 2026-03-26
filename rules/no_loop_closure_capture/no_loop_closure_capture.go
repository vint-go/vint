package no_loop_closure_capture

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoLoopClosureCaptureRule checks for references to enclosing loop variables
// from within nested functions (closures).
type NoLoopClosureCaptureRule struct{}

// Apply applies the rule to given file.
func (r *NoLoopClosureCaptureRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	// Since Go 1.22, loop variables are scoped per-iteration, so closure capture is no longer a bug.
	if file.Pkg.IsAtLeastGoVersion(lint.Go122) {
		return nil
	}

	var failures []lint.Failure

	w := &lintLoopClosureCapture{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoLoopClosureCaptureRule) Name() string {
	return "noLoopClosureCapture"
}

// Group returns the rule group.
func (*NoLoopClosureCaptureRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoLoopClosureCaptureRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintLoopClosureCapture struct {
	onFailure func(lint.Failure)
}

func (w *lintLoopClosureCapture) Visit(node ast.Node) ast.Visitor {
	// Collect loop variables based on loop type.
	var loopVars []*ast.Ident
	var body *ast.BlockStmt

	switch n := node.(type) {
	case *ast.RangeStmt:
		body = n.Body
		if n.Key != nil {
			if id, ok := n.Key.(*ast.Ident); ok {
				loopVars = append(loopVars, id)
			}
		}
		if n.Value != nil {
			if id, ok := n.Value.(*ast.Ident); ok {
				loopVars = append(loopVars, id)
			}
		}
	case *ast.ForStmt:
		body = n.Body
		switch post := n.Post.(type) {
		case *ast.AssignStmt:
			for _, lhs := range post.Lhs {
				if id, ok := lhs.(*ast.Ident); ok {
					loopVars = append(loopVars, id)
				}
			}
		case *ast.IncDecStmt:
			if id, ok := post.X.(*ast.Ident); ok {
				loopVars = append(loopVars, id)
			}
		}
	}

	if len(loopVars) == 0 || body == nil {
		return w
	}

	// Walk the loop body looking for closures that capture loop variables.
	w.checkBody(body, loopVars)

	// Don't recurse into this node's children via the default walker,
	// since checkBody already handles nested inspection.
	return nil
}

// checkBody inspects the loop body for closures that capture any of the given loop variables.
func (w *lintLoopClosureCapture) checkBody(body *ast.BlockStmt, loopVars []*ast.Ident) {
	// Track variables that are shadowed within the loop body.
	// A variable shadowed by := inside the loop body before any closure
	// should not be flagged.
	shadowedObjs := map[*ast.Object]bool{}

	for _, stmt := range body.List {
		// Check for shadowing assignments like v := v
		w.collectShadowed(stmt, loopVars, shadowedObjs)

		// Find closures (function literals) in this statement and check them.
		w.inspectForClosures(stmt, loopVars, shadowedObjs)
	}
}

// collectShadowed checks if a statement shadows any loop variables via short variable declaration.
func (w *lintLoopClosureCapture) collectShadowed(stmt ast.Stmt, loopVars []*ast.Ident, shadowed map[*ast.Object]bool) {
	assignStmt, ok := stmt.(*ast.AssignStmt)
	if !ok || assignStmt.Tok.String() != ":=" {
		return
	}
	for _, lhs := range assignStmt.Lhs {
		id, ok := lhs.(*ast.Ident)
		if !ok {
			continue
		}
		for _, v := range loopVars {
			if id.Name == v.Name && id.Obj != v.Obj {
				// This is a new variable that shadows the loop variable.
				shadowed[v.Obj] = true
			}
		}
	}
}

// inspectForClosures walks a statement tree to find function literals and check
// whether they reference any of the loop variables.
func (w *lintLoopClosureCapture) inspectForClosures(node ast.Node, loopVars []*ast.Ident, shadowed map[*ast.Object]bool) {
	ast.Inspect(node, func(n ast.Node) bool {
		// Look for go statements and defer statements with function literals.
		var callExpr *ast.CallExpr
		switch s := n.(type) {
		case *ast.GoStmt:
			callExpr = s.Call
		case *ast.DeferStmt:
			callExpr = s.Call
		default:
			return true
		}

		funcLit, ok := callExpr.Fun.(*ast.FuncLit)
		if !ok {
			return false
		}

		// Collect the parameter names of the function literal.
		// Variables passed as parameters are not captured.
		paramNames := map[string]bool{}
		if funcLit.Type != nil && funcLit.Type.Params != nil {
			for _, field := range funcLit.Type.Params.List {
				for _, name := range field.Names {
					paramNames[name.Name] = true
				}
			}
		}

		// Check the function literal body for references to loop variables.
		w.checkClosureBody(funcLit.Body, loopVars, shadowed, paramNames)

		return false
	})
}

// checkClosureBody checks a closure body for references to captured loop variables.
func (w *lintLoopClosureCapture) checkClosureBody(body *ast.BlockStmt, loopVars []*ast.Ident, shadowed map[*ast.Object]bool, paramNames map[string]bool) {
	if body == nil {
		return
	}

	ast.Inspect(body, func(n ast.Node) bool {
		// Skip key-value expressions (don't check keys).
		if kv, ok := n.(*ast.KeyValueExpr); ok {
			ast.Inspect(kv.Value, func(inner ast.Node) bool {
				id, ok := inner.(*ast.Ident)
				if !ok || id.Obj == nil {
					return true
				}
				w.checkIdent(id, loopVars, shadowed, paramNames)
				return true
			})
			return false
		}

		id, ok := n.(*ast.Ident)
		if !ok || id.Obj == nil {
			return true
		}

		w.checkIdent(id, loopVars, shadowed, paramNames)
		return true
	})
}

// checkIdent checks if an identifier references a captured loop variable.
func (w *lintLoopClosureCapture) checkIdent(id *ast.Ident, loopVars []*ast.Ident, shadowed map[*ast.Object]bool, paramNames map[string]bool) {
	// If the identifier name matches a parameter name, it's not captured.
	if paramNames[id.Name] {
		return
	}

	for _, v := range loopVars {
		if v.Obj == id.Obj && !shadowed[v.Obj] {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       id,
				Category:   lint.FailureCategoryLogic,
				Failure:    fmt.Sprintf("loop variable %s captured by closure", id.Name),
			})
			break
		}
	}
}

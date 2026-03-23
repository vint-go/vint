package no_unsafe_cors_bypass

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnsafeCorsBypassRule detects unsafe CORS bypass patterns such as reflecting
// the request Origin header without validation or using a wildcard
// Access-Control-Allow-Origin value.
type NoUnsafeCorsBypassRule struct{}

// Apply applies the rule to the given file.
func (r *NoUnsafeCorsBypassRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Quick check: the file must import "net/http" to be relevant.
	if !importsNetHTTP(file.AST) {
		return nil
	}

	onFailure := func(node ast.Node, msg string) {
		failures = append(failures, lint.Failure{
			Confidence: 1,
			Node:       node,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    msg,
		})
	}

	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		// Track variables that hold the Origin header value.
		originTainted := map[string]bool{}

		// Walk statements to find tainted variables and CORS header sets.
		analyzeBlock(fn.Body, originTainted, false, onFailure)
	}

	return failures
}

// Name returns the rule name.
func (*NoUnsafeCorsBypassRule) Name() string {
	return "noUnsafeCorsBypass"
}

// Group returns the rule group.
func (*NoUnsafeCorsBypassRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnsafeCorsBypassRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// importsNetHTTP checks if the file imports "net/http".
func importsNetHTTP(file *ast.File) bool {
	for _, imp := range file.Imports {
		if imp.Path.Value == `"net/http"` {
			return true
		}
	}
	return false
}

// analyzeBlock walks a block statement tracking Origin-tainted variables and
// checking for unsafe CORS header sets.
// The validated flag indicates that we are inside a conditional that validates
// the origin (e.g., if allowedOrigins[origin]).
func analyzeBlock(block *ast.BlockStmt, tainted map[string]bool, validated bool, onFailure func(ast.Node, string)) {
	if block == nil {
		return
	}
	for _, stmt := range block.List {
		analyzeStmt(stmt, tainted, validated, onFailure)
	}
}

// analyzeStmt processes a single statement.
func analyzeStmt(stmt ast.Stmt, tainted map[string]bool, validated bool, onFailure func(ast.Node, string)) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		// Check RHS for r.Header.Get("Origin") patterns.
		for i, rhs := range s.Rhs {
			if isOriginHeaderGet(rhs) && i < len(s.Lhs) {
				if ident, ok := s.Lhs[i].(*ast.Ident); ok && ident.Name != "_" {
					tainted[ident.Name] = true
				}
			}
		}
	case *ast.ExprStmt:
		// Check for w.Header().Set("Access-Control-Allow-Origin", ...)
		checkCORSHeaderSet(s.X, tainted, validated, onFailure)
	case *ast.IfStmt:
		if s.Init != nil {
			analyzeStmt(s.Init, tainted, validated, onFailure)
		}
		// If the condition references a tainted origin variable, the body
		// is considered validated (the developer is checking the origin).
		bodyValidated := validated || conditionReferencesTainted(s.Cond, tainted)
		analyzeBlock(s.Body, tainted, bodyValidated, onFailure)
		if s.Else != nil {
			// The else branch is NOT considered validated, since the check failed.
			if elseBlock, ok := s.Else.(*ast.BlockStmt); ok {
				analyzeBlock(elseBlock, tainted, validated, onFailure)
			} else if elseIf, ok := s.Else.(*ast.IfStmt); ok {
				analyzeStmt(elseIf, tainted, validated, onFailure)
			}
		}
	case *ast.ForStmt:
		if s.Init != nil {
			analyzeStmt(s.Init, tainted, validated, onFailure)
		}
		analyzeBlock(s.Body, tainted, validated, onFailure)
	case *ast.RangeStmt:
		analyzeBlock(s.Body, tainted, validated, onFailure)
	case *ast.BlockStmt:
		analyzeBlock(s, tainted, validated, onFailure)
	case *ast.SwitchStmt:
		analyzeBlock(s.Body, tainted, validated, onFailure)
	}
}

// conditionReferencesTainted checks if a condition expression references
// any tainted (origin) variable, suggesting the developer is validating it.
func conditionReferencesTainted(cond ast.Expr, tainted map[string]bool) bool {
	if cond == nil {
		return false
	}

	found := false
	ast.Inspect(cond, func(node ast.Node) bool {
		if found {
			return false
		}
		if ident, ok := node.(*ast.Ident); ok && tainted[ident.Name] {
			found = true
			return false
		}
		return true
	})
	return found
}

// isOriginHeaderGet checks if an expression matches patterns like:
//
//	r.Header.Get("Origin")
//
// where r is any variable.
func isOriginHeaderGet(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Get" {
		return false
	}

	// Check that the receiver is something.Header
	innerSel, ok := sel.X.(*ast.SelectorExpr)
	if !ok || innerSel.Sel.Name != "Header" {
		return false
	}

	// Check the argument is "Origin"
	if len(call.Args) != 1 {
		return false
	}
	return isStringLit(call.Args[0], "Origin")
}

// checkCORSHeaderSet checks if a call expression is setting the
// Access-Control-Allow-Origin header unsafely.
func checkCORSHeaderSet(expr ast.Expr, tainted map[string]bool, validated bool, onFailure func(ast.Node, string)) {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Set" || len(call.Args) != 2 {
		return
	}

	// Check first argument is "Access-Control-Allow-Origin"
	if !isStringLit(call.Args[0], "Access-Control-Allow-Origin") {
		return
	}

	// Check that the receiver is something.Header() -- i.e. w.Header().Set(...)
	innerCall, ok := sel.X.(*ast.CallExpr)
	if !ok {
		return
	}
	innerSel, ok := innerCall.Fun.(*ast.SelectorExpr)
	if !ok || innerSel.Sel.Name != "Header" {
		return
	}

	// Now check the second argument for unsafe patterns.
	value := call.Args[1]

	// Pattern 1: wildcard "*"
	if isStringLit(value, "*") {
		onFailure(call, "unsafe CORS bypass: Access-Control-Allow-Origin set to wildcard \"*\"")
		return
	}

	// Pattern 2: reflected origin (tainted variable) without validation
	if !validated && isOriginTainted(value, tainted) {
		onFailure(call, "unsafe CORS bypass: Access-Control-Allow-Origin reflects the request Origin header without validation")
		return
	}
}

// isOriginTainted checks if an expression is tainted with the Origin header value.
func isOriginTainted(expr ast.Expr, tainted map[string]bool) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		return tainted[e.Name]
	case *ast.CallExpr:
		// Direct inline: r.Header.Get("Origin")
		return isOriginHeaderGet(expr)
	case *ast.BinaryExpr:
		return isOriginTainted(e.X, tainted) || isOriginTainted(e.Y, tainted)
	case *ast.ParenExpr:
		return isOriginTainted(e.X, tainted)
	}
	return false
}

// isStringLit checks if the given expression is a string literal with the given value.
func isStringLit(expr ast.Expr, value string) bool {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return false
	}
	// String literals include the quotes, e.g., `"Origin"`.
	return lit.Value == `"`+value+`"`
}

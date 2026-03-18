package no_overwritten_argument

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoOverwrittenArgumentRule detects function arguments that are overwritten
// before they are ever read.
type NoOverwrittenArgumentRule struct{}

// Apply applies the rule to given file.
func (r *NoOverwrittenArgumentRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintOverwrittenArg{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoOverwrittenArgumentRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintOverwrittenArg{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoOverwrittenArgumentRule) Name() string {
	return "noOverwrittenArgument"
}

// Group returns the rule group.
func (*NoOverwrittenArgumentRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoOverwrittenArgumentRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintOverwrittenArg struct {
	onFailure func(lint.Failure)
}

func (w *lintOverwrittenArg) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Body != nil && n.Type.Params != nil {
			w.checkFunc(n.Type.Params, n.Body)
		}
		return nil
	case *ast.FuncLit:
		if n.Body != nil && n.Type.Params != nil {
			w.checkFunc(n.Type.Params, n.Body)
		}
		return nil
	}
	return w
}

// checkFunc checks if any parameter is overwritten before its first use.
func (w *lintOverwrittenArg) checkFunc(params *ast.FieldList, body *ast.BlockStmt) {
	// Collect parameter names.
	paramNames := map[string]bool{}
	for _, field := range params.List {
		for _, name := range field.Names {
			if name.Name != "_" {
				paramNames[name.Name] = true
			}
		}
	}

	if len(paramNames) == 0 {
		return
	}

	// Walk through top-level statements in the function body to find
	// parameters that are assigned before any read.
	// We track which params have been "read" and which have been "overwritten".
	read := map[string]bool{}
	overwritten := map[string]ast.Node{}

	for _, stmt := range body.List {
		w.analyzeStmt(stmt, paramNames, read, overwritten)
	}

	// Report failures for parameters overwritten before read.
	for name, node := range overwritten {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       node,
			Category:   lint.FailureCategoryLogic,
			Failure:    "argument '" + name + "' is overwritten before first use",
		})
	}
}

// analyzeStmt analyzes a statement to track reads and writes to parameter names.
// It processes statements sequentially; once a parameter is read, it is removed
// from further consideration. If a parameter is assigned before being read, it's
// recorded as overwritten.
func (w *lintOverwrittenArg) analyzeStmt(stmt ast.Stmt, params map[string]bool, read map[string]bool, overwritten map[string]ast.Node) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		// First check RHS for reads (RHS is evaluated before LHS assignment).
		for _, expr := range s.Rhs {
			w.markReads(expr, params, read)
		}
		// Then check LHS for writes.
		for _, expr := range s.Lhs {
			ident, ok := expr.(*ast.Ident)
			if !ok {
				// Could be a selector, index, etc. — mark reads.
				w.markReads(expr, params, read)
				continue
			}
			if !params[ident.Name] {
				continue
			}
			if read[ident.Name] {
				continue // already read, no problem
			}
			if _, already := overwritten[ident.Name]; already {
				continue // already reported
			}
			// For ":=", this is a new variable shadowing the param — not an overwrite of the param.
			if s.Tok == token.DEFINE {
				// Mark as read to prevent false positives on subsequent assignments
				// to the new variable (which shadows the param).
				read[ident.Name] = true
				continue
			}
			overwritten[ident.Name] = s
		}

	case *ast.ExprStmt:
		w.markReads(s.X, params, read)

	case *ast.ReturnStmt:
		for _, expr := range s.Results {
			w.markReads(expr, params, read)
		}

	case *ast.IfStmt:
		// If/else branches might conditionally read or write. Since the
		// parameter might be read in the condition or in a branch, we
		// conservatively mark reads in all parts.
		if s.Init != nil {
			w.markReadsInStmt(s.Init, params, read)
		}
		w.markReads(s.Cond, params, read)
		w.markReadsInBlock(s.Body, params, read)
		if s.Else != nil {
			w.markReadsInStmt(s.Else, params, read)
		}

	case *ast.SwitchStmt:
		if s.Init != nil {
			w.markReadsInStmt(s.Init, params, read)
		}
		if s.Tag != nil {
			w.markReads(s.Tag, params, read)
		}
		w.markReadsInBlock(s.Body, params, read)

	case *ast.TypeSwitchStmt:
		if s.Init != nil {
			w.markReadsInStmt(s.Init, params, read)
		}
		if s.Assign != nil {
			w.markReadsInStmt(s.Assign, params, read)
		}
		w.markReadsInBlock(s.Body, params, read)

	case *ast.ForStmt:
		if s.Init != nil {
			w.markReadsInStmt(s.Init, params, read)
		}
		if s.Cond != nil {
			w.markReads(s.Cond, params, read)
		}
		if s.Post != nil {
			w.markReadsInStmt(s.Post, params, read)
		}
		w.markReadsInBlock(s.Body, params, read)

	case *ast.RangeStmt:
		w.markReads(s.X, params, read)
		if s.Key != nil {
			w.markReads(s.Key, params, read)
		}
		if s.Value != nil {
			w.markReads(s.Value, params, read)
		}
		w.markReadsInBlock(s.Body, params, read)

	case *ast.GoStmt:
		w.markReads(s.Call, params, read)

	case *ast.DeferStmt:
		w.markReads(s.Call, params, read)

	case *ast.SendStmt:
		w.markReads(s.Chan, params, read)
		w.markReads(s.Value, params, read)

	case *ast.IncDecStmt:
		// i++ or i-- counts as both read and write.
		w.markReads(s.X, params, read)

	case *ast.DeclStmt:
		w.markReadsInDecl(s.Decl, params, read)

	case *ast.LabeledStmt:
		w.analyzeStmt(s.Stmt, params, read, overwritten)

	case *ast.BlockStmt:
		w.markReadsInBlock(s, params, read)

	case *ast.SelectStmt:
		w.markReadsInBlock(s.Body, params, read)

	case *ast.CaseClause:
		for _, expr := range s.List {
			w.markReads(expr, params, read)
		}
		for _, bodyStmt := range s.Body {
			w.markReadsInStmt(bodyStmt, params, read)
		}

	case *ast.CommClause:
		if s.Comm != nil {
			w.markReadsInStmt(s.Comm, params, read)
		}
		for _, bodyStmt := range s.Body {
			w.markReadsInStmt(bodyStmt, params, read)
		}
	}
}

// markReadsInBlock marks all parameter reads found inside a block statement.
func (w *lintOverwrittenArg) markReadsInBlock(block *ast.BlockStmt, params map[string]bool, read map[string]bool) {
	if block == nil {
		return
	}
	for _, stmt := range block.List {
		w.markReadsInStmt(stmt, params, read)
	}
}

// markReadsInStmt marks reads in any statement (delegates to markReads for expressions).
func (w *lintOverwrittenArg) markReadsInStmt(stmt ast.Stmt, params map[string]bool, read map[string]bool) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		for _, expr := range s.Rhs {
			w.markReads(expr, params, read)
		}
		for _, expr := range s.Lhs {
			w.markReads(expr, params, read)
		}
	case *ast.ExprStmt:
		w.markReads(s.X, params, read)
	case *ast.ReturnStmt:
		for _, expr := range s.Results {
			w.markReads(expr, params, read)
		}
	case *ast.IfStmt:
		if s.Init != nil {
			w.markReadsInStmt(s.Init, params, read)
		}
		w.markReads(s.Cond, params, read)
		w.markReadsInBlock(s.Body, params, read)
		if s.Else != nil {
			w.markReadsInStmt(s.Else, params, read)
		}
	case *ast.SwitchStmt:
		if s.Init != nil {
			w.markReadsInStmt(s.Init, params, read)
		}
		if s.Tag != nil {
			w.markReads(s.Tag, params, read)
		}
		w.markReadsInBlock(s.Body, params, read)
	case *ast.TypeSwitchStmt:
		if s.Init != nil {
			w.markReadsInStmt(s.Init, params, read)
		}
		if s.Assign != nil {
			w.markReadsInStmt(s.Assign, params, read)
		}
		w.markReadsInBlock(s.Body, params, read)
	case *ast.ForStmt:
		if s.Init != nil {
			w.markReadsInStmt(s.Init, params, read)
		}
		if s.Cond != nil {
			w.markReads(s.Cond, params, read)
		}
		if s.Post != nil {
			w.markReadsInStmt(s.Post, params, read)
		}
		w.markReadsInBlock(s.Body, params, read)
	case *ast.RangeStmt:
		w.markReads(s.X, params, read)
		if s.Key != nil {
			w.markReads(s.Key, params, read)
		}
		if s.Value != nil {
			w.markReads(s.Value, params, read)
		}
		w.markReadsInBlock(s.Body, params, read)
	case *ast.GoStmt:
		w.markReads(s.Call, params, read)
	case *ast.DeferStmt:
		w.markReads(s.Call, params, read)
	case *ast.SendStmt:
		w.markReads(s.Chan, params, read)
		w.markReads(s.Value, params, read)
	case *ast.IncDecStmt:
		w.markReads(s.X, params, read)
	case *ast.DeclStmt:
		w.markReadsInDecl(s.Decl, params, read)
	case *ast.LabeledStmt:
		w.markReadsInStmt(s.Stmt, params, read)
	case *ast.BlockStmt:
		w.markReadsInBlock(s, params, read)
	case *ast.SelectStmt:
		w.markReadsInBlock(s.Body, params, read)
	case *ast.CaseClause:
		for _, expr := range s.List {
			w.markReads(expr, params, read)
		}
		for _, bodyStmt := range s.Body {
			w.markReadsInStmt(bodyStmt, params, read)
		}
	case *ast.CommClause:
		if s.Comm != nil {
			w.markReadsInStmt(s.Comm, params, read)
		}
		for _, bodyStmt := range s.Body {
			w.markReadsInStmt(bodyStmt, params, read)
		}
	}
}

// markReadsInDecl marks reads in declaration statements.
func (w *lintOverwrittenArg) markReadsInDecl(decl ast.Decl, params map[string]bool, read map[string]bool) {
	gd, ok := decl.(*ast.GenDecl)
	if !ok {
		return
	}
	for _, spec := range gd.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for _, val := range vs.Values {
			w.markReads(val, params, read)
		}
	}
}

// markReads traverses an expression, marking any parameter identifiers as read.
func (w *lintOverwrittenArg) markReads(expr ast.Expr, params map[string]bool, read map[string]bool) {
	ast.Inspect(expr, func(n ast.Node) bool {
		// Don't descend into function literals — they have their own scope
		// and may capture the param but read it lazily.
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		if params[ident.Name] {
			read[ident.Name] = true
		}
		return true
	})
}

package no_sql_stmt_query

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSqlStmtQueryRule disallows calling (*database/sql.Stmt).Query without a context.
// The method (*sql.Stmt).Query does not accept a context.Context parameter,
// which prevents the caller from controlling cancellation, deadlines, and tracing.
// Use (*database/sql.Stmt).QueryContext instead, which accepts a context.Context as its first argument.
type NoSqlStmtQueryRule struct{}

// Apply applies the rule to given file.
func (r *NoSqlStmtQueryRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoSqlStmtQuery{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSqlStmtQueryRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	file.Pkg.TypeCheck()

	w := &lintNoSqlStmtQuery{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoSqlStmtQueryRule) Name() string {
	return "noSqlStmtQuery"
}

// Group returns the rule group.
func (*NoSqlStmtQueryRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoSqlStmtQueryRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoSqlStmtQueryRule) RequiresTypecheck() bool {
	return true
}

type lintNoSqlStmtQuery struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoSqlStmtQuery) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if this is a selector expression call: <expr>.Query(...)
	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	if sel.Sel.Name != "Query" {
		return w
	}

	// Check if the receiver type is *database/sql.Stmt
	t := w.pkg.TypeOf(sel.X)
	if t == nil {
		return w
	}

	if !isSqlStmtType(t) {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "(*sql.Stmt).Query does not accept a context; use (*sql.Stmt).QueryContext instead",
	})

	return w
}

// isSqlStmtType checks if the given type is *database/sql.Stmt.
func isSqlStmtType(t types.Type) bool {
	ptr, ok := t.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Name() == "Stmt" && obj.Pkg() != nil && obj.Pkg().Path() == "database/sql"
}

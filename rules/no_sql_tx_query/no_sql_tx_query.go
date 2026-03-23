package no_sql_tx_query

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSqlTxQueryRule disallows calling (*database/sql.Tx).Query without a context.
// The method (*sql.Tx).Query does not accept a context.Context parameter,
// which prevents the caller from controlling cancellation, deadlines, and tracing.
// Use (*database/sql.Tx).QueryContext instead, which accepts a context.Context as its first argument.
type NoSqlTxQueryRule struct{}

// Apply applies the rule to given file.
func (r *NoSqlTxQueryRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoSqlTxQuery{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSqlTxQueryRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	file.Pkg.TypeCheck()

	w := &lintNoSqlTxQuery{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoSqlTxQueryRule) Name() string {
	return "noSqlTxQuery"
}

// Group returns the rule group.
func (*NoSqlTxQueryRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoSqlTxQueryRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoSqlTxQueryRule) RequiresTypecheck() bool {
	return true
}

type lintNoSqlTxQuery struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoSqlTxQuery) Visit(node ast.Node) ast.Visitor {
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

	// Check if the receiver type is *database/sql.Tx
	t := w.pkg.TypeOf(sel.X)
	if t == nil {
		return w
	}

	if !isSqlTxType(t) {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "(*sql.Tx).Query does not accept a context; use (*sql.Tx).QueryContext instead",
	})

	return w
}

// isSqlTxType checks if the given type is *database/sql.Tx.
func isSqlTxType(t types.Type) bool {
	ptr, ok := t.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Name() == "Tx" && obj.Pkg() != nil && obj.Pkg().Path() == "database/sql"
}

package no_sql_tx_query_row

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSqlTxQueryRowRule disallows calling (*database/sql.Tx).QueryRow without a context.
// The method (*sql.Tx).QueryRow does not accept a context.Context parameter,
// which prevents the caller from controlling cancellation, deadlines, and tracing.
// Use (*database/sql.Tx).QueryRowContext instead, which accepts a context.Context as its first argument.
type NoSqlTxQueryRowRule struct{}

// Apply applies the rule to given file.
func (r *NoSqlTxQueryRowRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoSqlTxQueryRow{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSqlTxQueryRowRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	file.Pkg.TypeCheck()

	w := &lintNoSqlTxQueryRow{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoSqlTxQueryRowRule) Name() string {
	return "noSqlTxQueryRow"
}

// Group returns the rule group.
func (*NoSqlTxQueryRowRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoSqlTxQueryRowRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoSqlTxQueryRowRule) RequiresTypecheck() bool {
	return true
}

type lintNoSqlTxQueryRow struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoSqlTxQueryRow) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if this is a selector expression call: <expr>.QueryRow(...)
	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	if sel.Sel.Name != "QueryRow" {
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
		Failure:    "(*sql.Tx).QueryRow does not accept a context; use (*sql.Tx).QueryRowContext instead",
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

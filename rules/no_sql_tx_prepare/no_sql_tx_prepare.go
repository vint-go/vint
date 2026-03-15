package no_sql_tx_prepare

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoSqlTxPrepareRule disallows calling (*database/sql.Tx).Prepare without a context.
// The method (*sql.Tx).Prepare does not accept a context.Context parameter,
// which prevents the caller from controlling cancellation, deadlines, and tracing.
// Use (*database/sql.Tx).PrepareContext instead, which accepts a context.Context as its first argument.
type NoSqlTxPrepareRule struct{}

// Apply applies the rule to given file.
func (r *NoSqlTxPrepareRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoSqlTxPrepare{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSqlTxPrepareRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	file.Pkg.TypeCheck()

	w := &lintNoSqlTxPrepare{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoSqlTxPrepareRule) Name() string {
	return "noSqlTxPrepare"
}

// Group returns the rule group.
func (*NoSqlTxPrepareRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoSqlTxPrepareRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoSqlTxPrepareRule) RequiresTypecheck() bool {
	return true
}

type lintNoSqlTxPrepare struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoSqlTxPrepare) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if this is a selector expression call: <expr>.Prepare(...)
	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	if sel.Sel.Name != "Prepare" {
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
		Failure:    "(*sql.Tx).Prepare does not accept a context; use (*sql.Tx).PrepareContext instead",
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

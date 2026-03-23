package no_sql_db_prepare

import (
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSqlDbPrepareRule disallows calling (*database/sql.DB).Prepare without a context.
// The method (*sql.DB).Prepare does not accept a context.Context parameter,
// which prevents the caller from controlling cancellation, deadlines, and tracing.
// Use (*sql.DB).PrepareContext instead, which accepts a context.Context as its first argument.
type NoSqlDbPrepareRule struct{}

// Apply applies the rule to given file.
func (r *NoSqlDbPrepareRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoSqlDbPrepare{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSqlDbPrepareRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	file.Pkg.TypeCheck()

	w := &lintNoSqlDbPrepare{
		pkg:       file.Pkg,
		onFailure: onFailure,
	}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoSqlDbPrepareRule) Name() string {
	return "noSqlDbPrepare"
}

// Group returns the rule group.
func (*NoSqlDbPrepareRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoSqlDbPrepareRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoSqlDbPrepareRule) RequiresTypecheck() bool {
	return true
}

type lintNoSqlDbPrepare struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoSqlDbPrepare) Visit(node ast.Node) ast.Visitor {
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

	// Check if the receiver type is *database/sql.DB
	t := w.pkg.TypeOf(sel.X)
	if t == nil {
		return w
	}

	if !isSqlDBType(t) {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "(*sql.DB).Prepare does not accept a context; use (*sql.DB).PrepareContext instead",
	})

	return w
}

// isSqlDBType checks if the given type is *database/sql.DB.
func isSqlDBType(t types.Type) bool {
	ptr, ok := t.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}
	obj := named.Obj()
	return obj.Name() == "DB" && obj.Pkg() != nil && obj.Pkg().Path() == "database/sql"
}

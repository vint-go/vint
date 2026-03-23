package use_time_equal

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseTimeEqualRule flags where "==" and "!=" are used for equality checks on time.Time.
// time.Time values should be compared using the Equal method rather than ==,
// because == also compares the location, which may differ even for the same instant in time.
type UseTimeEqualRule struct{}

// Apply applies the rule to given file.
func (*UseTimeEqualRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.TypeCheck() != nil {
		return nil
	}

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintTimeEqual{file: file, onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseTimeEqualRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintTimeEqual{file: file, onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseTimeEqualRule) Name() string {
	return "useTimeEqual"
}

// Group returns the rule group.
func (*UseTimeEqualRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*UseTimeEqualRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule needs type info.
func (*UseTimeEqualRule) RequiresTypecheck() bool {
	return true
}

type lintTimeEqual struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (l *lintTimeEqual) Visit(node ast.Node) ast.Visitor {
	expr, ok := node.(*ast.BinaryExpr)
	if !ok {
		return l
	}

	switch expr.Op {
	case token.EQL, token.NEQ:
	default:
		return l
	}

	typeOfX := l.file.Pkg.TypeOf(expr.X)
	typeOfY := l.file.Pkg.TypeOf(expr.Y)
	if !isNamedType(typeOfX, "time", "Time") || !isNamedType(typeOfY, "time", "Time") {
		return l
	}

	negateStr := ""
	if expr.Op == token.NEQ {
		negateStr = "!"
	}

	l.onFailure(lint.Failure{
		Category:   lint.FailureCategoryTime,
		Confidence: 1,
		Node:       node,
		Failure:    fmt.Sprintf("use %s%s.Equal(%s) instead of %q operator", negateStr, astutils.GoFmt(expr.X), astutils.GoFmt(expr.Y), expr.Op),
	})

	return l
}

// isNamedType returns true if the given type is the named type from the given package.
func isNamedType(typ types.Type, importPath, name string) bool {
	if typ == nil {
		return false
	}
	n, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	obj := n.Obj()
	return obj != nil && obj.Pkg() != nil && obj.Pkg().Path() == importPath && obj.Name() == name
}

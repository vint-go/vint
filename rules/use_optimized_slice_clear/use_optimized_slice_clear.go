package use_optimized_slice_clear

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseOptimizedSliceClearRule detects slice clear loops that can use compiler-optimized idiom.
// Go's compiler can optimize `for i := range s { s[i] = zero }` patterns into a memclr call.
// This checker identifies loops that clear slices using `for i := 0; i < len(s); i++`
// and suggests using the optimized range-based pattern.
type UseOptimizedSliceClearRule struct{}

// Apply applies the rule to given file.
func (r *UseOptimizedSliceClearRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintOptimizedSliceClear{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseOptimizedSliceClearRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintOptimizedSliceClear{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseOptimizedSliceClearRule) Name() string {
	return "useOptimizedSliceClear"
}

// Group returns the rule group.
func (*UseOptimizedSliceClearRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*UseOptimizedSliceClearRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintOptimizedSliceClear struct {
	onFailure func(lint.Failure)
}

func (w *lintOptimizedSliceClear) Visit(node ast.Node) ast.Visitor {
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		return w
	}

	// Check pattern: for i := 0; i < len(s); i++ { s[i] = zero }
	sliceName := matchSliceClearLoop(forStmt)
	if sliceName == "" {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       forStmt,
		Category:   lint.FailureCategoryOptimization,
		Failure:    "use 'for i := range " + sliceName + "' to enable compiler memclr optimization",
	})

	return w
}

// matchSliceClearLoop checks if the for statement matches the pattern:
//
//	for i := 0; i < len(s); i++ { s[i] = zeroValue }
//
// Returns the slice name if matched, empty string otherwise.
func matchSliceClearLoop(forStmt *ast.ForStmt) string {
	// 1. Check Init: must be `i := 0`
	initAssign, ok := forStmt.Init.(*ast.AssignStmt)
	if !ok {
		return ""
	}
	if initAssign.Tok != token.DEFINE {
		return ""
	}
	if len(initAssign.Lhs) != 1 || len(initAssign.Rhs) != 1 {
		return ""
	}
	iterIdent, ok := initAssign.Lhs[0].(*ast.Ident)
	if !ok {
		return ""
	}
	// RHS must be 0
	rhsLit, ok := initAssign.Rhs[0].(*ast.BasicLit)
	if !ok || rhsLit.Kind != token.INT || rhsLit.Value != "0" {
		return ""
	}

	// 2. Check Cond: must be `i < len(s)`
	condBin, ok := forStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		return ""
	}
	if condBin.Op != token.LSS {
		return ""
	}
	// LHS of condition must be the same iterator
	condIdent, ok := condBin.X.(*ast.Ident)
	if !ok || condIdent.Name != iterIdent.Name {
		return ""
	}
	// RHS must be len(s)
	lenCall, ok := condBin.Y.(*ast.CallExpr)
	if !ok {
		return ""
	}
	lenIdent, ok := lenCall.Fun.(*ast.Ident)
	if !ok || lenIdent.Name != "len" {
		return ""
	}
	if len(lenCall.Args) != 1 {
		return ""
	}
	sliceName := astutils.GoFmt(lenCall.Args[0])
	if sliceName == "" {
		return ""
	}

	// 3. Check Post: must be `i++`
	postInc, ok := forStmt.Post.(*ast.IncDecStmt)
	if !ok {
		return ""
	}
	if postInc.Tok != token.INC {
		return ""
	}
	postIdent, ok := postInc.X.(*ast.Ident)
	if !ok || postIdent.Name != iterIdent.Name {
		return ""
	}

	// 4. Check Body: must have exactly one statement: s[i] = zeroValue
	if len(forStmt.Body.List) != 1 {
		return ""
	}
	bodyAssign, ok := forStmt.Body.List[0].(*ast.AssignStmt)
	if !ok {
		return ""
	}
	if bodyAssign.Tok != token.ASSIGN {
		return ""
	}
	if len(bodyAssign.Lhs) != 1 || len(bodyAssign.Rhs) != 1 {
		return ""
	}

	// LHS must be s[i]
	indexExpr, ok := bodyAssign.Lhs[0].(*ast.IndexExpr)
	if !ok {
		return ""
	}
	indexedSliceName := astutils.GoFmt(indexExpr.X)
	if indexedSliceName != sliceName {
		return ""
	}
	indexIdent, ok := indexExpr.Index.(*ast.Ident)
	if !ok || indexIdent.Name != iterIdent.Name {
		return ""
	}

	// RHS must be a zero value (0, "", nil, false, or zero struct literal)
	if !isZeroValue(bodyAssign.Rhs[0]) {
		return ""
	}

	return sliceName
}

// isZeroValue returns true if the expression represents a zero value.
func isZeroValue(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		// 0, 0.0, ""
		switch e.Kind {
		case token.INT:
			return e.Value == "0"
		case token.FLOAT:
			return e.Value == "0.0" || e.Value == "0."
		case token.STRING:
			return e.Value == `""`
		}
	case *ast.Ident:
		// nil, false
		return e.Name == "nil" || e.Name == "false"
	case *ast.CompositeLit:
		// Zero struct literal like T{}
		return len(e.Elts) == 0
	}
	return false
}

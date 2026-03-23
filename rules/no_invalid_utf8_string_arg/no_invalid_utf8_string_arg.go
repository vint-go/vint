package no_invalid_utf8_string_arg

import (
	"go/ast"
	"go/token"
	"strconv"
	"unicode/utf8"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoInvalidUtf8StringArgRule detects calls to strings package functions
// where invalid UTF-8 data is passed as an argument.
type NoInvalidUtf8StringArgRule struct{}

// Apply applies the rule to given file.
func (r *NoInvalidUtf8StringArgRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintInvalidUtf8{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoInvalidUtf8StringArgRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}
	w := &lintInvalidUtf8{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoInvalidUtf8StringArgRule) Name() string {
	return "noInvalidUtf8StringArg"
}

// Group returns the rule group.
func (*NoInvalidUtf8StringArgRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInvalidUtf8StringArgRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// stringsFuncsExpectingUTF8 is the set of strings package functions
// that expect valid UTF-8 input.
var stringsFuncsExpectingUTF8 = map[string]bool{
	"Compare":      true,
	"Contains":     true,
	"ContainsAny":  true,
	"ContainsRune": true,
	"Count":        true,
	"EqualFold":    true,
	"Fields":       true,
	"FieldsFunc":   true,
	"HasPrefix":    true,
	"HasSuffix":    true,
	"Index":        true,
	"IndexAny":     true,
	"IndexByte":    true,
	"IndexFunc":    true,
	"IndexRune":    true,
	"LastIndex":    true,
	"LastIndexAny": true,
	"LastIndexByte": true,
	"LastIndexFunc": true,
	"Map":          true,
	"Repeat":       true,
	"Replace":      true,
	"ReplaceAll":   true,
	"Split":        true,
	"SplitAfter":   true,
	"SplitAfterN":  true,
	"SplitN":       true,
	"Title":        true,
	"ToLower":      true,
	"ToLowerSpecial": true,
	"ToTitle":      true,
	"ToTitleSpecial": true,
	"ToUpper":      true,
	"ToUpperSpecial": true,
	"ToValidUTF8":  true,
	"Trim":         true,
	"TrimFunc":     true,
	"TrimLeft":     true,
	"TrimLeftFunc": true,
	"TrimPrefix":   true,
	"TrimRight":    true,
	"TrimRightFunc": true,
	"TrimSpace":    true,
	"TrimSuffix":   true,
}

type lintInvalidUtf8 struct {
	onFailure func(lint.Failure)
}

func (w *lintInvalidUtf8) Visit(node ast.Node) ast.Visitor {
	// We look at function bodies to track variable assignments.
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Body != nil {
			w.checkFuncBody(n.Body)
		}
		return nil
	case *ast.FuncLit:
		if n.Body != nil {
			w.checkFuncBody(n.Body)
		}
		return nil
	}
	return w
}

// checkFuncBody analyzes a function body for invalid UTF-8 strings passed to strings functions.
func (w *lintInvalidUtf8) checkFuncBody(body *ast.BlockStmt) {
	// Track variables that hold invalid UTF-8 strings.
	invalidVars := map[string]bool{}

	// Walk through statements collecting variable info and checking calls.
	for _, stmt := range body.List {
		w.analyzeStmt(stmt, invalidVars)
	}
}

// analyzeStmt recursively analyzes a statement.
func (w *lintInvalidUtf8) analyzeStmt(stmt ast.Stmt, invalidVars map[string]bool) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		// Track assignments of invalid UTF-8 strings.
		for i, rhs := range s.Rhs {
			if i < len(s.Lhs) {
				if ident, ok := s.Lhs[i].(*ast.Ident); ok {
					if isInvalidUTF8StringConversion(rhs) {
						invalidVars[ident.Name] = true
					}
				}
			}
		}
		// Also check if the RHS contains calls to strings functions with invalid args.
		w.checkExprsForCalls(s.Rhs, invalidVars)
	case *ast.ExprStmt:
		w.checkExprForCalls(s.X, invalidVars)
	case *ast.ReturnStmt:
		w.checkExprsForCalls(s.Results, invalidVars)
	case *ast.IfStmt:
		if s.Init != nil {
			w.analyzeStmt(s.Init, invalidVars)
		}
		w.checkExprForCalls(s.Cond, invalidVars)
		if s.Body != nil {
			for _, st := range s.Body.List {
				w.analyzeStmt(st, invalidVars)
			}
		}
		if s.Else != nil {
			w.analyzeStmt(s.Else, invalidVars)
		}
	case *ast.BlockStmt:
		for _, st := range s.List {
			w.analyzeStmt(st, invalidVars)
		}
	case *ast.ForStmt:
		if s.Init != nil {
			w.analyzeStmt(s.Init, invalidVars)
		}
		if s.Post != nil {
			w.analyzeStmt(s.Post, invalidVars)
		}
		if s.Body != nil {
			for _, st := range s.Body.List {
				w.analyzeStmt(st, invalidVars)
			}
		}
	case *ast.RangeStmt:
		if s.Body != nil {
			for _, st := range s.Body.List {
				w.analyzeStmt(st, invalidVars)
			}
		}
	case *ast.SwitchStmt:
		if s.Init != nil {
			w.analyzeStmt(s.Init, invalidVars)
		}
		if s.Body != nil {
			for _, st := range s.Body.List {
				w.analyzeStmt(st, invalidVars)
			}
		}
	case *ast.CaseClause:
		for _, st := range s.Body {
			w.analyzeStmt(st, invalidVars)
		}
	case *ast.DeclStmt:
		if gd, ok := s.Decl.(*ast.GenDecl); ok {
			for _, spec := range gd.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					for i, val := range vs.Values {
						if i < len(vs.Names) {
							if isInvalidUTF8StringConversion(val) {
								invalidVars[vs.Names[i].Name] = true
							}
						}
					}
					w.checkExprsForCalls(vs.Values, invalidVars)
				}
			}
		}
	case *ast.DeferStmt:
		w.checkExprForCalls(s.Call, invalidVars)
	case *ast.GoStmt:
		w.checkExprForCalls(s.Call, invalidVars)
	}
}

// checkExprsForCalls checks a slice of expressions for calls to strings functions with invalid UTF-8.
func (w *lintInvalidUtf8) checkExprsForCalls(exprs []ast.Expr, invalidVars map[string]bool) {
	for _, expr := range exprs {
		w.checkExprForCalls(expr, invalidVars)
	}
}

// checkExprForCalls recursively checks an expression for calls to strings functions with invalid UTF-8.
func (w *lintInvalidUtf8) checkExprForCalls(expr ast.Expr, invalidVars map[string]bool) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *ast.CallExpr:
		w.checkCallExpr(e, invalidVars)
		// Also check arguments recursively.
		for _, arg := range e.Args {
			w.checkExprForCalls(arg, invalidVars)
		}
		w.checkExprForCalls(e.Fun, invalidVars)
	case *ast.BinaryExpr:
		w.checkExprForCalls(e.X, invalidVars)
		w.checkExprForCalls(e.Y, invalidVars)
	case *ast.UnaryExpr:
		w.checkExprForCalls(e.X, invalidVars)
	case *ast.ParenExpr:
		w.checkExprForCalls(e.X, invalidVars)
	case *ast.IndexExpr:
		w.checkExprForCalls(e.X, invalidVars)
		w.checkExprForCalls(e.Index, invalidVars)
	}
}

// checkCallExpr checks if a call expression is a strings function call with invalid UTF-8 arguments.
func (w *lintInvalidUtf8) checkCallExpr(call *ast.CallExpr, invalidVars map[string]bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return
	}

	if pkgIdent.Name != "strings" {
		return
	}

	funcName := sel.Sel.Name
	if !stringsFuncsExpectingUTF8[funcName] {
		return
	}

	// Check each argument for invalid UTF-8.
	for _, arg := range call.Args {
		if w.isInvalidUTF8Arg(arg, invalidVars) {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       call,
				Failure:    "argument to strings." + funcName + " is not a valid UTF-8 encoded string",
			})
			return // one failure per call is enough
		}
	}
}

// isInvalidUTF8Arg checks if an expression is or refers to an invalid UTF-8 string.
func (w *lintInvalidUtf8) isInvalidUTF8Arg(expr ast.Expr, invalidVars map[string]bool) bool {
	// Direct string([]byte{...}) conversion.
	if isInvalidUTF8StringConversion(expr) {
		return true
	}

	// Variable that was assigned an invalid UTF-8 string.
	if ident, ok := expr.(*ast.Ident); ok {
		return invalidVars[ident.Name]
	}

	return false
}

// isInvalidUTF8StringConversion checks if an expression is a string([]byte{...}) conversion
// containing invalid UTF-8 bytes.
func isInvalidUTF8StringConversion(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	// Check it's a string() conversion.
	if !astutils.IsIdent(call.Fun, "string") {
		return false
	}

	if len(call.Args) != 1 {
		return false
	}

	// Check the argument is a []byte{...} composite literal.
	compLit, ok := call.Args[0].(*ast.CompositeLit)
	if !ok {
		return false
	}

	// Check the type is []byte.
	arrayType, ok := compLit.Type.(*ast.ArrayType)
	if !ok {
		return false
	}
	if !astutils.IsIdent(arrayType.Elt, "byte") {
		return false
	}

	// Extract byte values.
	bytes := make([]byte, 0, len(compLit.Elts))
	for _, elt := range compLit.Elts {
		lit, ok := elt.(*ast.BasicLit)
		if !ok {
			return false // can't statically determine
		}
		if lit.Kind != token.INT {
			return false
		}
		val, err := strconv.ParseUint(lit.Value, 0, 8)
		if err != nil {
			return false
		}
		bytes = append(bytes, byte(val))
	}

	if len(bytes) == 0 {
		return false
	}

	return !utf8.Valid(bytes)
}

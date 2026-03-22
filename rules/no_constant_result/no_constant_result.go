package no_constant_result

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoConstantResultRule reports named result parameters that always return
// the same constant value across all execution paths. By default, only
// unexported (private) functions are checked. Set check-exported to true
// to also analyze exported functions.
type NoConstantResultRule struct {
	checkExported bool
}

// Configure validates and applies the rule configuration.
func (r *NoConstantResultRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.checkExported = false
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noConstantResult" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		if isRuleOption(k, "check-exported") {
			val, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for check-exported in "noConstantResult" rule; need bool but got %T`, v)
			}
			r.checkExported = val
		}
	}

	return nil
}

// resultInfo tracks a single named result parameter of a function declaration.
type resultInfo struct {
	field     *ast.Field
	nameIdent *ast.Ident
	name      string
	index     int // positional index among all individual result names
}

func (r *NoConstantResultRule) analyzeFunc(funcDecl *ast.FuncDecl) []lint.Failure {
	results := collectNamedResults(funcDecl)
	if len(results) == 0 {
		return nil
	}

	returns := collectAllReturns(funcDecl.Body)
	if len(returns) == 0 {
		return nil
	}

	var failures []lint.Failure
	totalResults := countTotalResults(funcDecl)
	for _, ri := range results {
		constVal, isConst := analyzeResultConstancy(funcDecl.Body, ri, totalResults)
		if !isConst {
			continue
		}

		failures = append(failures, lint.Failure{
			Confidence: 1,
			Node:       funcDecl,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    fmt.Sprintf("result %s is always %s", ri.name, constVal),
		})
	}

	return failures
}

// Apply applies the rule to the given file.
func (r *NoConstantResultRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}

		if !r.checkExported && funcDecl.Name != nil && ast.IsExported(funcDecl.Name.Name) {
			continue
		}

		failures = append(failures, r.analyzeFunc(funcDecl)...)
	}

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoConstantResultRule) ApplyToNode(_ *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok || funcDecl.Body == nil {
		return nil
	}

	if !r.checkExported && funcDecl.Name != nil && ast.IsExported(funcDecl.Name.Name) {
		return nil
	}

	return r.analyzeFunc(funcDecl)
}

// Name returns the rule name.
func (*NoConstantResultRule) Name() string {
	return "noConstantResult"
}

// Group returns the rule group.
func (*NoConstantResultRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoConstantResultRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// analyzeResultConstancy determines if a named result is always returned as
// the same constant across all return paths, including bare returns and
// returns via the result's own name.
func analyzeResultConstancy(body *ast.BlockStmt, ri resultInfo, totalResults int) (string, bool) {
	returns := collectAllReturns(body)
	if len(returns) == 0 {
		return "", false
	}

	//nolint:staticcheck // ast.Object is deprecated but no replacement in go/ast
	tracker := trackResultAssignments(body, ri.nameIdent.Obj)

	// When tracker is tainted (non-constant modifications), we can still
	// detect explicit constant returns. But bare returns and named-result
	// returns are unreliable.
	if tracker.tainted {
		return checkExplicitReturnsOnly(returns, ri.index, totalResults)
	}

	zeroVal := zeroValueLiteral(ri.field)

	var seenValue string
	first := true

	for _, ret := range returns {
		var resolved string
		var ok bool

		if len(ret.Results) == 0 {
			// Bare return — resolve conservatively.
			resolved, ok = resolveForBareReturn(tracker, zeroVal)
		} else if len(ret.Results) != totalResults {
			return "", false
		} else {
			expr := ret.Results[ri.index]
			if isConstantExpr(expr) {
				resolved = astutils.GoFmt(expr)
				ok = true
			} else if ident, isIdent := expr.(*ast.Ident); isIdent && ident.Obj == ri.nameIdent.Obj { //nolint:staticcheck
				// Returning the named result by name.
				resolved, ok = tracker.constantValue()
			} else {
				return "", false
			}
		}

		if !ok {
			return "", false
		}

		if first {
			seenValue = resolved
			first = false
		} else if resolved != seenValue {
			return "", false
		}
	}

	if first {
		return "", false
	}
	return seenValue, true
}

// resolveForBareReturn determines what value a bare return produces for a
// named result. It is conservative: for non-zero tracked values, it gives
// up because unassigned code paths would return the zero value instead.
func resolveForBareReturn(tracker *assignTracker, zeroVal string) (string, bool) {
	if len(tracker.values) == 0 {
		// No assignments → zero value.
		if zeroVal != "" {
			return zeroVal, true
		}
		return "", false
	}
	// Has assignments → only safe if tracked constant equals zero value,
	// because paths that skip the assignment still return the zero value.
	constVal, ok := tracker.constantValue()
	if !ok {
		return "", false
	}
	if zeroVal != "" && constVal == zeroVal {
		return constVal, true
	}
	return "", false
}

// checkExplicitReturnsOnly falls back to checking only explicit constant
// returns (the original behavior) when assignment tracking is unreliable.
func checkExplicitReturnsOnly(returns []*ast.ReturnStmt, index, total int) (string, bool) {
	var seenValue string
	first := true

	for _, ret := range returns {
		if len(ret.Results) == 0 {
			// Bare return with tainted tracker → give up.
			return "", false
		}
		if len(ret.Results) != total {
			return "", false
		}
		expr := ret.Results[index]
		if !isConstantExpr(expr) {
			return "", false
		}
		rendered := astutils.GoFmt(expr)
		if first {
			seenValue = rendered
			first = false
		} else if rendered != seenValue {
			return "", false
		}
	}

	if first {
		return "", false
	}
	return seenValue, true
}

// assignTracker tracks assignments to a single named result variable.
type assignTracker struct {
	values  map[string]bool // set of distinct constant values assigned
	tainted bool            // true if any non-constant or compound modification
}

// constantValue returns the single constant value if all assignments
// are to the same constant.
func (t *assignTracker) constantValue() (string, bool) {
	if t.tainted || len(t.values) != 1 {
		return "", false
	}
	for v := range t.values {
		return v, true
	}
	return "", false
}

// trackResultAssignments walks a function body to find all modifications
// to a named result and classifies them. Skips nested function literals
// but taints if the result is captured by a closure.
//
//nolint:staticcheck // ast.Object is deprecated but no replacement in go/ast
func trackResultAssignments(body *ast.BlockStmt, resultObj *ast.Object) *assignTracker {
	tracker := &assignTracker{values: map[string]bool{}}

	// Pre-check: if the named result is captured by any closure, taint.
	if isCapturedByClosure(body, resultObj) {
		tracker.tainted = true
		return tracker
	}

	ast.Inspect(body, func(n ast.Node) bool {
		if tracker.tainted {
			return false
		}
		// Don't descend into nested function literals (already checked above).
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}

		switch node := n.(type) {
		case *ast.AssignStmt:
			for i, lhs := range node.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if !ok || ident.Obj != resultObj {
					continue
				}

				if node.Tok != token.ASSIGN {
					// +=, -=, etc.
					tracker.tainted = true
					return false
				}

				// Multi-value from function call (len(Lhs) != len(Rhs)).
				if len(node.Lhs) != len(node.Rhs) {
					tracker.tainted = true
					return false
				}

				rhs := node.Rhs[i]
				if isConstantExpr(rhs) {
					tracker.values[astutils.GoFmt(rhs)] = true
				} else {
					tracker.tainted = true
					return false
				}
			}

		case *ast.IncDecStmt:
			if ident, ok := node.X.(*ast.Ident); ok && ident.Obj == resultObj {
				tracker.tainted = true
				return false
			}

		case *ast.UnaryExpr:
			// Address-of: &result — could be modified externally.
			if node.Op == token.AND {
				if ident, ok := node.X.(*ast.Ident); ok && ident.Obj == resultObj {
					tracker.tainted = true
					return false
				}
			}
		}

		return true
	})

	return tracker
}

// isCapturedByClosure checks if a named result is referenced inside any
// closure (function literal) in the body.
//
//nolint:staticcheck // ast.Object is deprecated but no replacement in go/ast
func isCapturedByClosure(body *ast.BlockStmt, resultObj *ast.Object) bool {
	captured := false
	ast.Inspect(body, func(n ast.Node) bool {
		if captured {
			return false
		}
		fl, ok := n.(*ast.FuncLit)
		if !ok {
			return true
		}
		ast.Inspect(fl.Body, func(inner ast.Node) bool {
			if captured {
				return false
			}
			if ident, ok := inner.(*ast.Ident); ok && ident.Obj == resultObj {
				captured = true
			}
			return true
		})
		return false // don't re-enter the func lit in the outer walk
	})
	return captured
}

// collectNamedResults extracts a flat list of named result parameters
// from a function declaration.
func collectNamedResults(fn *ast.FuncDecl) []resultInfo {
	if fn.Type.Results == nil {
		return nil
	}
	var result []resultInfo
	idx := 0
	for _, field := range fn.Type.Results.List {
		if len(field.Names) == 0 {
			idx++
			continue
		}
		for _, name := range field.Names {
			if name.Name == "_" {
				idx++
				continue
			}
			result = append(result, resultInfo{
				field:     field,
				nameIdent: name,
				name:      name.Name,
				index:     idx,
			})
			idx++
		}
	}
	return result
}

// countTotalResults counts the total number of result values (named and
// unnamed) in a function's result list.
func countTotalResults(fn *ast.FuncDecl) int {
	if fn.Type.Results == nil {
		return 0
	}
	count := 0
	for _, field := range fn.Type.Results.List {
		if len(field.Names) == 0 {
			count++
		} else {
			count += len(field.Names)
		}
	}
	return count
}

// collectAllReturns collects all return statements in a function body
// (including bare returns), excluding those in nested function literals.
func collectAllReturns(body *ast.BlockStmt) []*ast.ReturnStmt {
	var stmts []*ast.ReturnStmt
	ast.Inspect(body, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}
		if ret, ok := n.(*ast.ReturnStmt); ok {
			stmts = append(stmts, ret)
		}
		return true
	})
	return stmts
}

// zeroValueLiteral returns the Go literal for the zero value of a type,
// or "" if the type is not recognized.
func zeroValueLiteral(field *ast.Field) string {
	switch t := field.Type.(type) {
	case *ast.Ident:
		switch t.Name {
		case "error":
			return "nil"
		case "bool":
			return "false"
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64", "complex64", "complex128",
			"byte", "rune", "uintptr":
			return "0"
		case "string":
			return `""`
		}
	case *ast.StarExpr:
		return "nil"
	case *ast.InterfaceType:
		return "nil"
	case *ast.MapType:
		return "nil"
	case *ast.ChanType:
		return "nil"
	case *ast.FuncType:
		return "nil"
	case *ast.ArrayType:
		if t.Len == nil {
			return "nil" // slice
		}
	}
	return ""
}

// isConstantExpr returns true if the expression is a compile-time constant
// literal (basic lit, unary of basic lit, identifiers like true/false/nil).
func isConstantExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return true
	case *ast.Ident:
		return e.Name == "true" || e.Name == "false" || e.Name == "nil"
	case *ast.UnaryExpr:
		return isConstantExpr(e.X)
	case *ast.ParenExpr:
		return isConstantExpr(e.X)
	default:
		return false
	}
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}

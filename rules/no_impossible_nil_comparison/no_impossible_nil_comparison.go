package no_impossible_nil_comparison

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoImpossibleNilComparisonRule detects comparisons of interface-typed values
// with nil that are always false (== nil) or always true (!= nil). This happens
// when the value was produced by a function that provably never returns untyped
// nil, or when a concrete type is implicitly wrapped in an interface.
//
// This matches the semantics of staticcheck SA4023: checking at the call site
// rather than at the return site.
type NoImpossibleNilComparisonRule struct{}

// Apply applies the rule to the given file.
func (r *NoImpossibleNilComparisonRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if file.Pkg.TypeCheck() != nil {
		return nil
	}

	info := file.Pkg.TypesInfo()
	if info == nil {
		return nil
	}

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	// Phase 1: Analyze all functions defined in this file to determine which
	// ones provably never return untyped nil at their interface return positions.
	neverNil := buildNeverNilMap(file, info)

	// Phase 2: Walk the AST checking nil comparisons at call sites.
	w := &nilCompChecker{
		file:      file,
		info:      info,
		neverNil:  neverNil,
		onFailure: onFailure,
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoImpossibleNilComparisonRule) Name() string {
	return "noImpossibleNilComparison"
}

// Group returns the rule group.
func (*NoImpossibleNilComparisonRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoImpossibleNilComparisonRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type info.
func (*NoImpossibleNilComparisonRule) RequiresTypecheck() bool {
	return true
}

// --- Phase 1: Build map of functions that never return nil ---

// neverNilInfo describes a function that never returns a nil interface value
// at certain return positions.
type neverNilInfo struct {
	positions map[int]bool // return positions that are provably never nil
	name      string       // function name for diagnostic messages
}

// buildNeverNilMap analyzes all function declarations in the file and identifies
// which ones never return untyped nil at their interface-typed return positions.
func buildNeverNilMap(file *lint.File, info *types.Info) map[types.Object]*neverNilInfo {
	result := make(map[types.Object]*neverNilInfo)

	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil || fn.Name == nil {
			continue
		}

		obj := info.ObjectOf(fn.Name)
		if obj == nil {
			continue
		}

		sig, ok := obj.Type().(*types.Signature)
		if !ok {
			continue
		}

		nni := analyzeFunc(sig, fn.Body, file)
		if nni != nil {
			nni.name = fn.Name.Name
			result[obj] = nni
		}
	}

	return result
}

// analyzeFunc checks whether a function never returns untyped nil at any of
// its interface-typed return positions.
func analyzeFunc(sig *types.Signature, body *ast.BlockStmt, file *lint.File) *neverNilInfo {
	results := sig.Results()
	if results == nil || results.Len() == 0 {
		return nil
	}

	// Find which return positions have interface types
	ifacePositions := make(map[int]bool)
	for i := 0; i < results.Len(); i++ {
		if _, isIface := results.At(i).Type().Underlying().(*types.Interface); isIface {
			ifacePositions[i] = true
		}
	}
	if len(ifacePositions) == 0 {
		return nil
	}

	// Collect all return statements (skip nested function literals)
	returns := collectReturns(body)
	if len(returns) == 0 {
		return nil
	}

	// For each interface position, check if any return could produce nil
	neverNilPos := make(map[int]bool)
	for pos := range ifacePositions {
		canBeNil := false
		for _, ret := range returns {
			if canReturnNilAt(ret, pos, file, results.Len()) {
				canBeNil = true
				break
			}
		}
		if !canBeNil {
			neverNilPos[pos] = true
		}
	}

	if len(neverNilPos) == 0 {
		return nil
	}

	return &neverNilInfo{positions: neverNilPos}
}

// collectReturns collects all return statements in a function body,
// excluding those inside nested function literals.
func collectReturns(body *ast.BlockStmt) []*ast.ReturnStmt {
	var returns []*ast.ReturnStmt
	ast.Inspect(body, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}
		if ret, ok := n.(*ast.ReturnStmt); ok {
			returns = append(returns, ret)
		}
		return true
	})
	return returns
}

// canReturnNilAt checks whether a return statement could produce a nil
// interface value at the given return position.
func canReturnNilAt(ret *ast.ReturnStmt, pos int, file *lint.File, resultCount int) bool {
	// Bare return (named results) - conservatively assume possible nil
	if len(ret.Results) == 0 {
		return true
	}

	// Single expression returning multiple values - conservatively assume possible nil
	if len(ret.Results) == 1 && resultCount > 1 {
		return true
	}

	if pos >= len(ret.Results) {
		return true
	}

	expr := ret.Results[pos]

	// Untyped nil produces a nil interface
	if isNilIdent(expr) {
		return true
	}

	exprType := file.Pkg.TypeOf(expr)
	if exprType == nil || exprType == types.Typ[types.Invalid] {
		return true // conservatively assume possible nil
	}

	// Interface-typed expression could be nil
	if _, isIface := exprType.Underlying().(*types.Interface); isIface {
		return true
	}

	// Concrete type wrapped in interface is always non-nil
	return false
}

// --- Phase 2: Check nil comparisons at call sites ---

// nilCompChecker walks the AST looking for impossible nil comparisons.
type nilCompChecker struct {
	file      *lint.File
	info      *types.Info
	neverNil  map[types.Object]*neverNilInfo
	onFailure func(lint.Failure)
}

func (w *nilCompChecker) Visit(node ast.Node) ast.Visitor {
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

// assignInfo tracks a variable's assignment from a function call.
type assignInfo struct {
	calledFunc types.Object     // the function/method that was called
	sig        *types.Signature // its signature
	position   int              // which return value position this variable got
	count      int              // number of assignments to this variable
}

// checkFuncBody analyzes a function body for impossible nil comparisons.
func (w *nilCompChecker) checkFuncBody(body *ast.BlockStmt) {
	// Step 1: Collect variable assignments from function calls
	assignments := w.collectAssignments(body)

	// Step 2: Find nil comparisons and check them
	ast.Inspect(body, func(n ast.Node) bool {
		// Process nested function literals recursively
		if fl, ok := n.(*ast.FuncLit); ok {
			if fl.Body != nil {
				w.checkFuncBody(fl.Body)
			}
			return false
		}

		binExpr, ok := n.(*ast.BinaryExpr)
		if !ok || (binExpr.Op != token.EQL && binExpr.Op != token.NEQ) {
			return true
		}

		// One side must be nil
		var nonNilExpr ast.Expr
		if isNilIdent(binExpr.X) {
			nonNilExpr = binExpr.Y
		} else if isNilIdent(binExpr.Y) {
			nonNilExpr = binExpr.X
		} else {
			return true
		}

		// The non-nil side must have an interface type
		nonNilType := w.file.Pkg.TypeOf(nonNilExpr)
		if nonNilType == nil || nonNilType == types.Typ[types.Invalid] {
			return true
		}
		if _, isIface := nonNilType.Underlying().(*types.Interface); !isIface {
			return true
		}

		alwaysResult := "false"
		if binExpr.Op == token.NEQ {
			alwaysResult = "true"
		}

		// Case 1: Direct function call: f() == nil
		if call, ok := nonNilExpr.(*ast.CallExpr); ok {
			w.checkDirectCall(binExpr, call, alwaysResult)
			return true
		}

		// Case 2: Variable from function call: x == nil
		if ident, ok := nonNilExpr.(*ast.Ident); ok {
			w.checkVarComparison(binExpr, ident, assignments, alwaysResult)
		}

		return true
	})
}

// collectAssignments scans a function body for variable assignments from
// function calls and builds a tracking map. Variables with multiple assignments
// are tracked but won't be flagged (conservative approach).
func (w *nilCompChecker) collectAssignments(body *ast.BlockStmt) map[types.Object]*assignInfo {
	result := make(map[types.Object]*assignInfo)

	ast.Inspect(body, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			return false // skip nested function literals
		}

		switch node := n.(type) {
		case *ast.AssignStmt:
			w.trackAssignment(node, result)
		case *ast.ValueSpec:
			w.trackValueSpec(node, result)
		}

		return true
	})

	return result
}

// trackAssignment tracks variable assignments from function calls in
// regular assignment statements (x := f() or x = f()).
func (w *nilCompChecker) trackAssignment(assign *ast.AssignStmt, result map[types.Object]*assignInfo) {
	// Only track single RHS function calls: x, y := f()
	if len(assign.Rhs) != 1 {
		w.bumpAssignCounts(assign.Lhs, result)
		return
	}

	call, isCall := assign.Rhs[0].(*ast.CallExpr)
	if !isCall {
		w.bumpAssignCounts(assign.Lhs, result)
		return
	}

	calledFunc := resolveCallTarget(call, w.info)
	var sig *types.Signature
	if calledFunc != nil {
		sig, _ = calledFunc.Type().(*types.Signature)
	}

	for i, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok || ident.Name == "_" {
			continue
		}
		obj := w.info.ObjectOf(ident)
		if obj == nil {
			continue
		}

		if existing, exists := result[obj]; exists {
			existing.count++
			existing.calledFunc = calledFunc
			existing.sig = sig
			existing.position = i
		} else {
			result[obj] = &assignInfo{
				calledFunc: calledFunc,
				sig:        sig,
				position:   i,
				count:      1,
			}
		}
	}
}

// trackValueSpec tracks variable assignments from function calls in
// var declarations (var x Type = f()).
func (w *nilCompChecker) trackValueSpec(vs *ast.ValueSpec, result map[types.Object]*assignInfo) {
	if len(vs.Values) != 1 {
		return
	}

	call, isCall := vs.Values[0].(*ast.CallExpr)
	if !isCall {
		return
	}

	calledFunc := resolveCallTarget(call, w.info)
	var sig *types.Signature
	if calledFunc != nil {
		sig, _ = calledFunc.Type().(*types.Signature)
	}

	for i, name := range vs.Names {
		if name.Name == "_" {
			continue
		}
		obj := w.info.ObjectOf(name)
		if obj == nil {
			continue
		}

		if existing, exists := result[obj]; exists {
			existing.count++
			existing.calledFunc = calledFunc
			existing.sig = sig
			existing.position = i
		} else {
			result[obj] = &assignInfo{
				calledFunc: calledFunc,
				sig:        sig,
				position:   i,
				count:      1,
			}
		}
	}
}

// bumpAssignCounts increments assignment counts for variables assigned
// non-call values, invalidating their tracking for the never-nil check.
func (w *nilCompChecker) bumpAssignCounts(lhsList []ast.Expr, result map[types.Object]*assignInfo) {
	for _, lhs := range lhsList {
		ident, ok := lhs.(*ast.Ident)
		if !ok || ident.Name == "_" {
			continue
		}
		obj := w.info.ObjectOf(ident)
		if obj == nil {
			continue
		}
		if existing, exists := result[obj]; exists {
			existing.count++
		} else {
			result[obj] = &assignInfo{count: 1}
		}
	}
}

// checkDirectCall checks if a direct function call comparison (f() == nil)
// is impossible.
func (w *nilCompChecker) checkDirectCall(binExpr *ast.BinaryExpr, call *ast.CallExpr, alwaysResult string) {
	calledFunc := resolveCallTarget(call, w.info)
	if calledFunc == nil {
		return
	}

	// Check if function is in our never-nil map (same-file body analysis)
	if nni, exists := w.neverNil[calledFunc]; exists && nni.positions[0] {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 0.8,
			Node:       binExpr,
			Failure:    fmt.Sprintf("nil comparison is always %s because %s never returns nil", alwaysResult, nni.name),
		})
		return
	}

	// Check if function's return type is concrete (produces non-nil interface)
	sig, ok := calledFunc.Type().(*types.Signature)
	if !ok {
		return
	}
	results := sig.Results()
	if results == nil || results.Len() == 0 {
		return
	}
	retType := results.At(0).Type()
	if retType == nil || retType == types.Typ[types.Invalid] {
		return
	}
	if _, isIface := retType.Underlying().(*types.Interface); !isIface {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 0.8,
			Node:       binExpr,
			Failure:    fmt.Sprintf("nil comparison is always %s because %s returns concrete type %s", alwaysResult, calledFunc.Name(), retType),
		})
	}
}

// checkVarComparison checks if a variable nil comparison (x == nil) is
// impossible because the variable was assigned from a function that never
// returns nil.
func (w *nilCompChecker) checkVarComparison(binExpr *ast.BinaryExpr, ident *ast.Ident, assignments map[types.Object]*assignInfo, alwaysResult string) {
	obj := w.info.ObjectOf(ident)
	if obj == nil {
		return
	}

	ai, exists := assignments[obj]
	if !exists || ai.count != 1 || ai.calledFunc == nil {
		return
	}

	// Check if function is in our never-nil map (same-file body analysis)
	if nni, nniExists := w.neverNil[ai.calledFunc]; nniExists && nni.positions[ai.position] {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 0.8,
			Node:       binExpr,
			Failure:    fmt.Sprintf("nil comparison of %s is always %s because %s never returns nil", ident.Name, alwaysResult, nni.name),
		})
		return
	}

	// Check if function's return type at this position is concrete
	if ai.sig == nil {
		return
	}
	results := ai.sig.Results()
	if results == nil || ai.position >= results.Len() {
		return
	}
	retType := results.At(ai.position).Type()
	if retType == nil || retType == types.Typ[types.Invalid] {
		return
	}
	if _, isIface := retType.Underlying().(*types.Interface); !isIface {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 0.8,
			Node:       binExpr,
			Failure:    fmt.Sprintf("nil comparison of %s is always %s because %s returns concrete type %s", ident.Name, alwaysResult, ai.calledFunc.Name(), retType),
		})
	}
}

// resolveCallTarget resolves a function call expression to its types.Object.
func resolveCallTarget(call *ast.CallExpr, info *types.Info) types.Object {
	switch fn := call.Fun.(type) {
	case *ast.Ident:
		return info.ObjectOf(fn)
	case *ast.SelectorExpr:
		return info.ObjectOf(fn.Sel)
	}
	return nil
}

// isNilIdent checks if an expression is the nil identifier.
func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}

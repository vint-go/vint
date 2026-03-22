package no_variable_shadowing

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoVariableShadowingRule checks for shadowed variables via short variable declarations.
type NoVariableShadowingRule struct {
	strict bool
}

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configure implements the [lint.ConfigurableRule] interface.
func (r *NoVariableShadowingRule) Configure(arguments lint.Arguments) error {
	// Default to strict mode (report all shadowing).
	r.strict = true
	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noVariableShadowing" rule, expecting a k,v map, got %T`, arguments[0])
	}

	if strictVal, ok := argKV["strict"]; ok {
		b, ok := strictVal.(bool)
		if !ok {
			return fmt.Errorf(`invalid "strict" value in "noVariableShadowing" rule; need bool but got %T`, strictVal)
		}
		r.strict = b
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoVariableShadowingRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	var pending []pendingShadow

	w := &lintVariableShadowing{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
		scopes:         nil,
		strict:         r.strict,
		pendingShadows: &pending,
	}

	// Collect top-level scope names from the file scope.
	if file.AST.Scope != nil {
		topScope := make(map[string]token.Pos)
		for name, obj := range file.AST.Scope.Objects {
			topScope[name] = obj.Pos()
		}
		w.scopes = append(w.scopes, topScope)
	}

	ast.Walk(w, file.AST)

	return failures
}

// Note: this rule specifically is a bad idea to implement with ApplyToNode,
// because of file.AST.Scope pass it would do down below for every Node..

// Name returns the rule name.
func (*NoVariableShadowingRule) Name() string {
	return "noVariableShadowing"
}

// Group returns the rule group.
func (*NoVariableShadowingRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoVariableShadowingRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// pendingShadow records a shadow that needs to be confirmed against
// remaining statements in non-strict mode.
type pendingShadow struct {
	failure    lint.Failure
	name       string
	scopeDepth int // depth (index) of the scope where the shadowed outer variable lives
}

type lintVariableShadowing struct {
	onFailure func(lint.Failure)
	// scopes is a stack of variable name -> declaration position maps.
	// Each entry represents an enclosing scope.
	scopes []map[string]token.Pos
	strict bool
	// pendingShadows collects shadows in non-strict mode for deferred resolution.
	// Shared across inner walkers via pointer so shadows bubble up.
	pendingShadows *[]pendingShadow
}

func (w *lintVariableShadowing) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		w.visitFunc(n.Type, n.Body)
		return nil // handled manually
	case *ast.FuncLit:
		w.visitFunc(n.Type, n.Body)
		return nil // handled manually
	}
	return w
}

// visitFunc handles entering a function scope. It collects parameters and named results
// into the same scope as the function body, since Go treats them as one scope.
func (w *lintVariableShadowing) visitFunc(funcType *ast.FuncType, body *ast.BlockStmt) {
	if body == nil {
		return
	}

	// In Go, function parameters, named return values, and the function body
	// all share the same scope. So we create a single scope that holds all of them.
	funcScope := make(map[string]token.Pos)

	// Collect parameters.
	if funcType.Params != nil {
		for _, field := range funcType.Params.List {
			for _, name := range field.Names {
				if name.Name != "_" {
					funcScope[name.Name] = name.Pos()
				}
			}
		}
	}

	// Collect named return values.
	if funcType.Results != nil {
		for _, field := range funcType.Results.List {
			for _, name := range field.Names {
				if name.Name != "_" {
					funcScope[name.Name] = name.Pos()
				}
			}
		}
	}

	// Create a new walker with the function scope pushed.
	// We walk the body statements directly (not via walkBlock) to avoid
	// creating an extra scope level, since params and body share a scope.
	inner := &lintVariableShadowing{
		onFailure:      w.onFailure,
		scopes:         append(copyScopes(w.scopes), funcScope),
		strict:         w.strict,
		pendingShadows: w.pendingShadows,
	}

	inner.walkStmtList(body.List)
}

// walkBlock walks a block statement, creating a new scope for it.
func (w *lintVariableShadowing) walkBlock(block *ast.BlockStmt) {
	if block == nil {
		return
	}

	// This block introduces a new scope.
	blockScope := make(map[string]token.Pos)
	inner := &lintVariableShadowing{
		onFailure:      w.onFailure,
		scopes:         append(copyScopes(w.scopes), blockScope),
		strict:         w.strict,
		pendingShadows: w.pendingShadows,
	}

	inner.walkStmtList(block.List)
}

// walkStmtList walks a list of statements sequentially. In non-strict mode,
// it resolves pending shadows by checking if the shadowed outer variable is
// referenced in subsequent statements. Shadows are only resolved at the scope
// level where the shadowed variable was defined — shadows of variables from
// deeper outer scopes pass through to be resolved at the appropriate level.
func (w *lintVariableShadowing) walkStmtList(stmts []ast.Stmt) {
	// The current scope depth is the index of the innermost scope.
	currentDepth := len(w.scopes) - 1

	for i, stmt := range stmts {
		beforeLen := len(*w.pendingShadows)
		w.walkStmt(stmt)

		if !w.strict && len(*w.pendingShadows) > beforeLen {
			remaining := stmts[i+1:]
			newPending := (*w.pendingShadows)[beforeLen:]
			var kept []pendingShadow
			for _, ps := range newPending {
				if ps.scopeDepth == currentDepth {
					// The shadowed variable lives in OUR scope — resolve here.
					if identUsedInStmts(ps.name, remaining) {
						w.onFailure(ps.failure)
					}
				} else {
					// The shadowed variable lives further out — pass through.
					kept = append(kept, ps)
				}
			}
			*w.pendingShadows = append((*w.pendingShadows)[:beforeLen], kept...)
		}
	}
}

// walkStmt walks a single statement, checking for shadowing in assignments
// and recursing into sub-blocks.
func (w *lintVariableShadowing) walkStmt(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if s.Tok == token.DEFINE { // :=
			w.checkShortVarDecl(s)
		}
		// Check RHS for function literals.
		for _, rhs := range s.Rhs {
			w.walkExprForFuncLit(rhs)
		}

	case *ast.DeclStmt:
		if genDecl, ok := s.Decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range vs.Names {
						if name.Name != "_" {
							w.addToCurrentScope(name.Name, name.Pos())
						}
					}
					// Check values for function literals.
					for _, val := range vs.Values {
						w.walkExprForFuncLit(val)
					}
				}
			}
		}

	case *ast.BlockStmt:
		w.walkBlock(s)

	case *ast.IfStmt:
		w.walkIfStmt(s)

	case *ast.ForStmt:
		w.walkForStmt(s)

	case *ast.RangeStmt:
		w.walkRangeStmt(s)

	case *ast.SwitchStmt:
		w.walkSwitchStmt(s)

	case *ast.TypeSwitchStmt:
		w.walkTypeSwitchStmt(s)

	case *ast.SelectStmt:
		if s.Body != nil {
			for _, cc := range s.Body.List {
				if comm, ok := cc.(*ast.CommClause); ok {
					w.walkCommClause(comm)
				}
			}
		}

	case *ast.GoStmt:
		w.walkExprForFuncLit(s.Call.Fun)

	case *ast.DeferStmt:
		w.walkExprForFuncLit(s.Call.Fun)

	case *ast.ExprStmt:
		w.walkExprForFuncLit(s.X)

	case *ast.ReturnStmt:
		for _, expr := range s.Results {
			w.walkExprForFuncLit(expr)
		}

	case *ast.SendStmt:
		w.walkExprForFuncLit(s.Value)
	}
}

// walkExprForFuncLit checks an expression for function literals and walks them.
func (w *lintVariableShadowing) walkExprForFuncLit(expr ast.Expr) {
	ast.Inspect(expr, func(n ast.Node) bool {
		if fl, ok := n.(*ast.FuncLit); ok {
			inner := &lintVariableShadowing{
				onFailure:      w.onFailure,
				scopes:         copyScopes(w.scopes),
				strict:         w.strict,
				pendingShadows: w.pendingShadows,
			}
			inner.visitFunc(fl.Type, fl.Body)
			return false
		}
		return true
	})
}

// walkIfStmt handles if-init statements. The init and body share a scope.
func (w *lintVariableShadowing) walkIfStmt(s *ast.IfStmt) {
	ifScope := make(map[string]token.Pos)
	inner := &lintVariableShadowing{
		onFailure:      w.onFailure,
		scopes:         append(copyScopes(w.scopes), ifScope),
		strict:         w.strict,
		pendingShadows: w.pendingShadows,
	}

	if s.Init != nil {
		inner.walkStmt(s.Init)
	}

	inner.walkBlock(s.Body)

	if s.Else != nil {
		inner.walkStmt(s.Else)
	}
}

// walkForStmt handles for statements.
func (w *lintVariableShadowing) walkForStmt(s *ast.ForStmt) {
	forScope := make(map[string]token.Pos)
	inner := &lintVariableShadowing{
		onFailure:      w.onFailure,
		scopes:         append(copyScopes(w.scopes), forScope),
		strict:         w.strict,
		pendingShadows: w.pendingShadows,
	}

	if s.Init != nil {
		inner.walkStmt(s.Init)
	}
	inner.walkBlock(s.Body)
}

// walkRangeStmt handles range statements.
func (w *lintVariableShadowing) walkRangeStmt(s *ast.RangeStmt) {
	rangeScope := make(map[string]token.Pos)
	inner := &lintVariableShadowing{
		onFailure:      w.onFailure,
		scopes:         append(copyScopes(w.scopes), rangeScope),
		strict:         w.strict,
		pendingShadows: w.pendingShadows,
	}

	if s.Tok == token.DEFINE {
		if key, ok := s.Key.(*ast.Ident); ok && key.Name != "_" {
			inner.checkShadowAndAdd(key)
		}
		if val, ok := s.Value.(*ast.Ident); ok && val.Name != "_" {
			inner.checkShadowAndAdd(val)
		}
	}

	inner.walkBlock(s.Body)
}

// walkSwitchStmt handles switch statements.
func (w *lintVariableShadowing) walkSwitchStmt(s *ast.SwitchStmt) {
	switchScope := make(map[string]token.Pos)
	inner := &lintVariableShadowing{
		onFailure:      w.onFailure,
		scopes:         append(copyScopes(w.scopes), switchScope),
		strict:         w.strict,
		pendingShadows: w.pendingShadows,
	}

	if s.Init != nil {
		inner.walkStmt(s.Init)
	}

	if s.Body != nil {
		for _, cc := range s.Body.List {
			if clause, ok := cc.(*ast.CaseClause); ok {
				inner.walkCaseClause(clause)
			}
		}
	}
}

// walkTypeSwitchStmt handles type switch statements.
func (w *lintVariableShadowing) walkTypeSwitchStmt(s *ast.TypeSwitchStmt) {
	switchScope := make(map[string]token.Pos)
	inner := &lintVariableShadowing{
		onFailure:      w.onFailure,
		scopes:         append(copyScopes(w.scopes), switchScope),
		strict:         w.strict,
		pendingShadows: w.pendingShadows,
	}

	if s.Init != nil {
		inner.walkStmt(s.Init)
	}

	if s.Assign != nil {
		inner.walkStmt(s.Assign)
	}

	if s.Body != nil {
		for _, cc := range s.Body.List {
			if clause, ok := cc.(*ast.CaseClause); ok {
				inner.walkCaseClause(clause)
			}
		}
	}
}

// walkCaseClause handles case clauses in switch statements.
func (w *lintVariableShadowing) walkCaseClause(clause *ast.CaseClause) {
	caseScope := make(map[string]token.Pos)
	inner := &lintVariableShadowing{
		onFailure:      w.onFailure,
		scopes:         append(copyScopes(w.scopes), caseScope),
		strict:         w.strict,
		pendingShadows: w.pendingShadows,
	}

	inner.walkStmtList(clause.Body)
}

// walkCommClause handles comm clauses in select statements.
func (w *lintVariableShadowing) walkCommClause(clause *ast.CommClause) {
	commScope := make(map[string]token.Pos)
	inner := &lintVariableShadowing{
		onFailure:      w.onFailure,
		scopes:         append(copyScopes(w.scopes), commScope),
		strict:         w.strict,
		pendingShadows: w.pendingShadows,
	}

	if clause.Comm != nil {
		inner.walkStmt(clause.Comm)
	}

	inner.walkStmtList(clause.Body)
}

// checkShortVarDecl checks a short variable declaration for shadowed variables.
func (w *lintVariableShadowing) checkShortVarDecl(assign *ast.AssignStmt) {
	for _, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok || ident.Name == "_" {
			continue
		}

		name := ident.Name

		// Check if this variable already exists in the CURRENT scope.
		// If it does, this := is just reusing the existing variable (Go allows
		// := with at least one new variable on the LHS, reusing others).
		if w.existsInCurrentScope(name) {
			continue
		}

		w.checkShadowAndAdd(ident)
	}
}

// checkShadowAndAdd checks if ident shadows a variable from an outer scope,
// reports it if so, and adds it to the current scope.
func (w *lintVariableShadowing) checkShadowAndAdd(ident *ast.Ident) {
	name := ident.Name

	// Check outer scopes (all scopes except the current/innermost one)
	// for a variable with the same name.
	if len(w.scopes) > 1 {
		for i := len(w.scopes) - 2; i >= 0; i-- {
			if _, exists := w.scopes[i][name]; exists {
				f := lint.Failure{
					Confidence: 1,
					Node:       ident,
					Category:   lint.FailureCategoryLogic,
					Failure:    fmt.Sprintf("variable %s shadows variable from outer scope", name),
				}
				if w.strict {
					w.onFailure(f)
				} else {
					*w.pendingShadows = append(*w.pendingShadows, pendingShadow{
						failure:    f,
						name:       name,
						scopeDepth: i,
					})
				}
				break
			}
		}
	}

	// Add to the current (innermost) scope.
	w.addToCurrentScope(name, ident.Pos())
}

// existsInCurrentScope checks if a name exists in the current (innermost) scope.
func (w *lintVariableShadowing) existsInCurrentScope(name string) bool {
	if len(w.scopes) == 0 {
		return false
	}
	_, exists := w.scopes[len(w.scopes)-1][name]
	return exists
}

// addToCurrentScope adds a variable to the innermost scope.
func (w *lintVariableShadowing) addToCurrentScope(name string, pos token.Pos) {
	if len(w.scopes) > 0 {
		w.scopes[len(w.scopes)-1][name] = pos
	}
}

// copyScopes creates a shallow copy of the scope stack.
func copyScopes(scopes []map[string]token.Pos) []map[string]token.Pos {
	if scopes == nil {
		return nil
	}
	result := make([]map[string]token.Pos, len(scopes))
	copy(result, scopes)
	return result
}

// identUsedInStmts checks whether a variable name is referenced as an identifier
// anywhere within the given statements. Used in non-strict mode to determine
// if a shadowed outer variable is still relevant after the shadowing scope.
func identUsedInStmts(name string, stmts []ast.Stmt) bool {
	for _, stmt := range stmts {
		if identUsedInNode(name, stmt) {
			return true
		}
	}
	return false
}

// identUsedInNode checks whether a variable name is referenced as an identifier
// anywhere within the given AST node.
func identUsedInNode(name string, node ast.Node) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if found {
			return false
		}
		if ident, ok := n.(*ast.Ident); ok && ident.Name == name {
			found = true
			return false
		}
		return true
	})
	return found
}

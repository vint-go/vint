package no_nil_map_assignment

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoNilMapAssignmentRule warns when assigning to a nil map, which causes a runtime panic.
type NoNilMapAssignmentRule struct{}

// Apply applies the rule to the given file.
func (r *NoNilMapAssignmentRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoNilMapAssignment{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		w.checkFunctionBody(funcDecl.Body)
	}

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoNilMapAssignmentRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok || funcDecl.Body == nil {
		return nil
	}

	var failures []lint.Failure
	w := &lintNoNilMapAssignment{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.checkFunctionBody(funcDecl.Body)
	return failures
}

// Name returns the rule name.
func (*NoNilMapAssignmentRule) Name() string {
	return "noNilMapAssignment"
}

// Group returns the rule group.
func (*NoNilMapAssignmentRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoNilMapAssignmentRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoNilMapAssignment struct {
	onFailure func(lint.Failure)
}

// checkFunctionBody examines a function body for assignments to nil maps.
func (w *lintNoNilMapAssignment) checkFunctionBody(body *ast.BlockStmt) {
	// Collect variables declared as nil maps via "var m map[...]..."
	nilMaps := w.collectNilMapVars(body.List)
	if len(nilMaps) == 0 {
		return
	}

	// Walk the body looking for index assignments to those variables
	w.checkStatements(body.List, nilMaps)
}

// collectNilMapVars finds all variable names declared with map types and no initial value
// in the given statement list (non-recursive for top-level declarations).
func (w *lintNoNilMapAssignment) collectNilMapVars(stmts []ast.Stmt) map[string]bool {
	nilMaps := map[string]bool{}

	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.DeclStmt:
			genDecl, ok := s.Decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range genDecl.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				// Check if it's a map type with no initialization values
				if _, isMap := vs.Type.(*ast.MapType); isMap && len(vs.Values) == 0 {
					for _, name := range vs.Names {
						nilMaps[name.Name] = true
					}
				}
			}
		}
	}

	return nilMaps
}

// checkStatements walks statements looking for index assignments to nil map variables.
// It also tracks when a variable gets reassigned (initialized), removing it from the nil set.
func (w *lintNoNilMapAssignment) checkStatements(stmts []ast.Stmt, nilMaps map[string]bool) {
	for _, stmt := range stmts {
		w.checkStmt(stmt, nilMaps)
	}
}

func (w *lintNoNilMapAssignment) checkStmt(stmt ast.Stmt, nilMaps map[string]bool) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		w.checkAssign(s, nilMaps)
	case *ast.BlockStmt:
		if s != nil {
			w.checkStatements(s.List, nilMaps)
		}
	case *ast.IfStmt:
		if s.Init != nil {
			w.checkStmt(s.Init, nilMaps)
		}
		if s.Body != nil {
			// Use a copy so that initializations inside if branches
			// don't affect the outer scope tracking
			w.checkStatements(s.Body.List, copyMap(nilMaps))
		}
		if s.Else != nil {
			w.checkStmt(s.Else, copyMap(nilMaps))
		}
	case *ast.ForStmt:
		if s.Init != nil {
			w.checkStmt(s.Init, nilMaps)
		}
		if s.Body != nil {
			w.checkStatements(s.Body.List, copyMap(nilMaps))
		}
	case *ast.RangeStmt:
		if s.Body != nil {
			w.checkStatements(s.Body.List, copyMap(nilMaps))
		}
	case *ast.SwitchStmt:
		if s.Init != nil {
			w.checkStmt(s.Init, nilMaps)
		}
		if s.Body != nil {
			for _, cc := range s.Body.List {
				clause, ok := cc.(*ast.CaseClause)
				if !ok {
					continue
				}
				w.checkStatements(clause.Body, copyMap(nilMaps))
			}
		}
	case *ast.TypeSwitchStmt:
		if s.Init != nil {
			w.checkStmt(s.Init, nilMaps)
		}
		if s.Body != nil {
			for _, cc := range s.Body.List {
				clause, ok := cc.(*ast.CaseClause)
				if !ok {
					continue
				}
				w.checkStatements(clause.Body, copyMap(nilMaps))
			}
		}
	case *ast.SelectStmt:
		if s.Body != nil {
			for _, cc := range s.Body.List {
				clause, ok := cc.(*ast.CommClause)
				if !ok {
					continue
				}
				w.checkStatements(clause.Body, copyMap(nilMaps))
			}
		}
	case *ast.DeclStmt:
		// Handle var declarations inside the function body (nested)
		genDecl, ok := s.Decl.(*ast.GenDecl)
		if !ok {
			return
		}
		for _, spec := range genDecl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			if _, isMap := vs.Type.(*ast.MapType); isMap && len(vs.Values) == 0 {
				for _, name := range vs.Names {
					nilMaps[name.Name] = true
				}
			}
		}
	}
}

// checkAssign checks an assignment statement for nil map index assignment and
// tracks variable re-initialization.
func (w *lintNoNilMapAssignment) checkAssign(assign *ast.AssignStmt, nilMaps map[string]bool) {
	// First, check if any LHS is an index expression into a nil map
	for _, lhs := range assign.Lhs {
		indexExpr, ok := lhs.(*ast.IndexExpr)
		if !ok {
			continue
		}
		ident, ok := indexExpr.X.(*ast.Ident)
		if !ok {
			continue
		}
		if nilMaps[ident.Name] {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       assign,
				Failure:    "assignment to nil map",
			})
		}
	}

	// Track reassignments that initialize the variable (make, composite literal, or assignment from another variable)
	for i, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok {
			continue
		}
		if !nilMaps[ident.Name] {
			continue
		}
		// If the variable is being reassigned (not via index), it's no longer nil
		if i < len(assign.Rhs) {
			delete(nilMaps, ident.Name)
		} else if len(assign.Rhs) == 1 {
			// Multiple LHS, single RHS (e.g., multi-return)
			delete(nilMaps, ident.Name)
		}
	}
}

func copyMap(m map[string]bool) map[string]bool {
	c := make(map[string]bool, len(m))
	for k, v := range m {
		c[k] = v
	}
	return c
}

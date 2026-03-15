package no_ineffectual_assignment

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoIneffectualAssignmentRule detects ineffectual assignments in Go code.
// An assignment is ineffectual if the variable assigned is not thereafter used
// before being reassigned or going out of scope.
type NoIneffectualAssignmentRule struct {
	checkEscapingErrors bool
}

// Configure validates and applies the rule configuration.
func (r *NoIneffectualAssignmentRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.checkEscapingErrors = false
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noIneffectualAssignment" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		if isRuleOption(k, "check-escaping-errors") {
			val, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for check-escaping-errors in "noIneffectualAssignment" rule; need bool but got %T`, v)
			}
			r.checkEscapingErrors = val
		}
	}

	return nil
}

// Apply applies the rule to the given file.
func (r *NoIneffectualAssignmentRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintIneffectualAssignment{
		checkEscapingErrors: r.checkEscapingErrors,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		w.checkFunc(fn)
	}

	return failures
}

// Name returns the rule name.
func (*NoIneffectualAssignmentRule) Name() string {
	return "noIneffectualAssignment"
}

// Group returns the rule group.
func (*NoIneffectualAssignmentRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoIneffectualAssignmentRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintIneffectualAssignment struct {
	checkEscapingErrors bool
	onFailure           func(lint.Failure)
}

// varState tracks whether a variable's latest assignment has been read.
type varState struct {
	assigned   bool     // whether the variable has been assigned a value
	read       bool     // whether the variable has been read since last assignment
	node       ast.Node // the assignment node (for reporting)
	name       string   // variable name
	fromParent bool     // whether this state was inherited from a parent scope (branch)
}

// scope tracks variable states within a block scope.
type scope struct {
	vars map[string]*varState
}

func newScope() *scope {
	return &scope{vars: make(map[string]*varState)}
}

// cloneForBranch creates a clone of this scope, marking all entries as fromParent.
func (s *scope) cloneForBranch() *scope {
	ns := newScope()
	for k, v := range s.vars {
		copied := *v
		copied.fromParent = true
		ns.vars[k] = &copied
	}
	return ns
}

func (w *lintIneffectualAssignment) checkFunc(fn *ast.FuncDecl) {
	if fn.Body == nil {
		return
	}

	sc := newScope()

	// Register function parameters as assigned+read (they come from the caller)
	if fn.Type.Params != nil {
		for _, field := range fn.Type.Params.List {
			for _, name := range field.Names {
				if name.Name != "_" {
					sc.vars[name.Name] = &varState{assigned: true, read: true, name: name.Name}
				}
			}
		}
	}

	// Register named return values as assigned+read (zero-initialized, used implicitly)
	if fn.Type.Results != nil {
		for _, field := range fn.Type.Results.List {
			for _, name := range field.Names {
				if name.Name != "_" {
					sc.vars[name.Name] = &varState{assigned: true, read: true, name: name.Name}
				}
			}
		}
	}

	w.analyzeBlock(fn.Body.List, sc, false)
}

func (w *lintIneffectualAssignment) analyzeBlock(stmts []ast.Stmt, sc *scope, inLoop bool) {
	for _, stmt := range stmts {
		w.analyzeStmt(stmt, sc, inLoop)
	}
}

func (w *lintIneffectualAssignment) analyzeStmt(stmt ast.Stmt, sc *scope, inLoop bool) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *ast.AssignStmt:
		// First, mark reads on the RHS
		for _, rhs := range s.Rhs {
			w.markReads(rhs, sc)
		}
		// Then process LHS assignments
		for _, lhs := range s.Lhs {
			ident, ok := lhs.(*ast.Ident)
			if !ok || ident.Name == "_" {
				// Selector expressions (struct fields) are not flagged
				// as stated in the spec
				w.markReads(lhs, sc)
				continue
			}
			// Check if there is a prior assignment that was not read
			if vs, exists := sc.vars[ident.Name]; exists && vs.assigned && !vs.read && !inLoop {
				// Only report if this is not an inherited parent-scope assignment
				// (parent-scope assignments in branches are handled by mergeScopes)
				if !vs.fromParent {
					w.reportFailure(vs)
				}
			}
			// Update state: new assignment, not yet read
			sc.vars[ident.Name] = &varState{
				assigned: true,
				read:     false,
				node:     s,
				name:     ident.Name,
			}
		}

	case *ast.DeclStmt:
		gd, ok := s.Decl.(*ast.GenDecl)
		if !ok {
			return
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			// Mark reads on values
			for _, val := range vs.Values {
				w.markReads(val, sc)
			}
			for _, name := range vs.Names {
				if name.Name == "_" {
					continue
				}
				// Variable declared with initial value
				if len(vs.Values) > 0 {
					sc.vars[name.Name] = &varState{
						assigned: true,
						read:     false,
						node:     vs,
						name:     name.Name,
					}
				} else {
					// var x int - zero value, not considered ineffectual
					sc.vars[name.Name] = &varState{
						assigned: false,
						read:     false,
						node:     vs,
						name:     name.Name,
					}
				}
			}
		}

	case *ast.IfStmt:
		if s.Init != nil {
			w.analyzeStmt(s.Init, sc, inLoop)
		}
		w.markReads(s.Cond, sc)

		thenScope := sc.cloneForBranch()
		w.analyzeBlock(s.Body.List, thenScope, inLoop)

		if s.Else != nil {
			elseScope := sc.cloneForBranch()
			w.analyzeStmt(s.Else, elseScope, inLoop)
			// After if/else, merge both branches
			w.mergeScopes(sc, thenScope, elseScope)
		} else {
			// Only then branch, merge conservatively
			w.mergeScopeSingle(sc, thenScope)
		}

	case *ast.BlockStmt:
		w.analyzeBlock(s.List, sc, inLoop)

	case *ast.ForStmt:
		if s.Init != nil {
			w.analyzeStmt(s.Init, sc, inLoop)
		}
		if s.Cond != nil {
			w.markReads(s.Cond, sc)
		}
		loopScope := sc.cloneForBranch()
		w.analyzeBlock(s.Body.List, loopScope, true)
		if s.Post != nil {
			w.analyzeStmt(s.Post, loopScope, true)
		}
		// After loop, merge conservatively (loop body may or may not execute)
		w.mergeScopeSingle(sc, loopScope)

	case *ast.RangeStmt:
		w.markReads(s.X, sc)
		// Range key/value are assignments
		if s.Key != nil {
			if ident, ok := s.Key.(*ast.Ident); ok && ident.Name != "_" {
				sc.vars[ident.Name] = &varState{assigned: true, read: false, node: s, name: ident.Name}
			}
		}
		if s.Value != nil {
			if ident, ok := s.Value.(*ast.Ident); ok && ident.Name != "_" {
				sc.vars[ident.Name] = &varState{assigned: true, read: false, node: s, name: ident.Name}
			}
		}
		loopScope := sc.cloneForBranch()
		w.analyzeBlock(s.Body.List, loopScope, true)
		w.mergeScopeSingle(sc, loopScope)

	case *ast.SwitchStmt:
		if s.Init != nil {
			w.analyzeStmt(s.Init, sc, inLoop)
		}
		if s.Tag != nil {
			w.markReads(s.Tag, sc)
		}
		w.analyzeSwitchBody(s.Body, sc, inLoop)

	case *ast.TypeSwitchStmt:
		if s.Init != nil {
			w.analyzeStmt(s.Init, sc, inLoop)
		}
		if s.Assign != nil {
			w.analyzeStmt(s.Assign, sc, inLoop)
		}
		w.analyzeSwitchBody(s.Body, sc, inLoop)

	case *ast.SelectStmt:
		if s.Body != nil {
			w.analyzeSwitchBody(s.Body, sc, inLoop)
		}

	case *ast.ReturnStmt:
		for _, result := range s.Results {
			w.markReads(result, sc)
		}

	case *ast.ExprStmt:
		w.markReads(s.X, sc)

	case *ast.SendStmt:
		w.markReads(s.Chan, sc)
		w.markReads(s.Value, sc)

	case *ast.IncDecStmt:
		// x++ or x-- reads and writes x
		w.markReads(s.X, sc)
		if ident, ok := s.X.(*ast.Ident); ok && ident.Name != "_" {
			sc.vars[ident.Name] = &varState{
				assigned: true,
				read:     false,
				node:     s,
				name:     ident.Name,
			}
		}

	case *ast.GoStmt:
		w.markReads(s.Call, sc)

	case *ast.DeferStmt:
		w.markReads(s.Call, sc)

	case *ast.LabeledStmt:
		w.analyzeStmt(s.Stmt, sc, inLoop)

	case *ast.BranchStmt:
		// break, continue, goto, fallthrough - no reads/writes to track

	case *ast.CommClause:
		if s.Comm != nil {
			w.analyzeStmt(s.Comm, sc, inLoop)
		}
		w.analyzeBlock(s.Body, sc, inLoop)

	case *ast.CaseClause:
		for _, expr := range s.List {
			w.markReads(expr, sc)
		}
		w.analyzeBlock(s.Body, sc, inLoop)
	}
}

func (w *lintIneffectualAssignment) analyzeSwitchBody(body *ast.BlockStmt, sc *scope, inLoop bool) {
	if body == nil {
		return
	}
	var branchScopes []*scope
	for _, stmt := range body.List {
		branchScope := sc.cloneForBranch()
		w.analyzeStmt(stmt, branchScope, inLoop)
		branchScopes = append(branchScopes, branchScope)
	}
	// Merge all branches conservatively
	for _, bs := range branchScopes {
		w.mergeScopeSingle(sc, bs)
	}
}

// mergeScopes merges two branch scopes (then/else) into the parent.
// If both branches assign a variable without reading it first, the parent's
// assignment was ineffectual. If either branch reads the variable, mark it read.
func (w *lintIneffectualAssignment) mergeScopes(parent, branch1, branch2 *scope) {
	for name, vs := range parent.vars {
		b1, has1 := branch1.vars[name]
		b2, has2 := branch2.vars[name]

		b1Read := has1 && b1.read
		b2Read := has2 && b2.read

		if b1Read || b2Read {
			vs.read = true
		} else if has1 && has2 && b1.assigned && b2.assigned && !b1.read && !b2.read {
			// Both branches assigned without reading -> parent assignment is dead
			// The parent's assignment is ineffectual; flag it if it was an actual assignment
			if vs.assigned && !vs.read {
				w.reportFailure(vs)
				// Mark as read so we don't report again
				vs.read = true
			}
		}
	}
	// Bring in new variables introduced in branches
	for name, b1vs := range branch1.vars {
		if _, exists := parent.vars[name]; !exists {
			parent.vars[name] = b1vs
		}
	}
	for name, b2vs := range branch2.vars {
		if _, exists := parent.vars[name]; !exists {
			parent.vars[name] = b2vs
		}
	}
}

// mergeScopeSingle merges a single branch scope back to the parent.
// Conservative: if the branch read a var, mark it read in parent.
func (w *lintIneffectualAssignment) mergeScopeSingle(parent, branch *scope) {
	for name, vs := range parent.vars {
		if bvs, exists := branch.vars[name]; exists && bvs.read {
			vs.read = true
		}
	}
	// Bring in new variables introduced in branch
	for name, bvs := range branch.vars {
		if _, exists := parent.vars[name]; !exists {
			parent.vars[name] = bvs
		}
	}
}

// markReads walks an expression and marks any referenced variables as read.
func (w *lintIneffectualAssignment) markReads(expr ast.Expr, sc *scope) {
	if expr == nil {
		return
	}

	ast.Inspect(expr, func(n ast.Node) bool {
		switch e := n.(type) {
		case *ast.Ident:
			if vs, exists := sc.vars[e.Name]; exists {
				vs.read = true
			}
		case *ast.FuncLit:
			// Variables captured by closures count as read
			w.markClosureCaptures(e, sc)
			return false
		}
		return true
	})
}

// markClosureCaptures marks all variables referenced in a function literal as read.
func (w *lintIneffectualAssignment) markClosureCaptures(fl *ast.FuncLit, sc *scope) {
	ast.Inspect(fl, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			if vs, exists := sc.vars[ident.Name]; exists {
				vs.read = true
			}
		}
		return true
	})
}

func (w *lintIneffectualAssignment) reportFailure(vs *varState) {
	w.onFailure(lint.Failure{
		Confidence: 1,
		Category:   lint.FailureCategoryLogic,
		Failure:    fmt.Sprintf("ineffectual assignment to %s", vs.name),
		Node:       vs.node,
	})
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}

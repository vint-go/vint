package no_high_cyclomatic_complexity

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/strowk/vint/lint"
)

// NoHighCyclomaticComplexityRule checks the cyclomatic complexity of Go functions
// and reports those that exceed a configurable threshold.
type NoHighCyclomaticComplexityRule struct {
	minComplexity int
}

const defaultMinCyclomaticComplexity = 30

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoHighCyclomaticComplexityRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.minComplexity = defaultMinCyclomaticComplexity
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		// try direct int64 argument
		complexity, ok := arguments[0].(int64)
		if !ok {
			return fmt.Errorf(`invalid argument to the "noHighCyclomaticComplexity" rule, expecting a k,v map or int64, got %T`, arguments[0])
		}
		r.minComplexity = int(complexity)
		return nil
	}

	r.minComplexity = defaultMinCyclomaticComplexity
	for k, v := range argKV {
		if isRuleOption(k, "min-complexity") {
			complexity, ok := v.(int64)
			if !ok {
				return fmt.Errorf(`invalid configuration value for min-complexity in "noHighCyclomaticComplexity" rule; need int64 but got %T`, v)
			}
			r.minComplexity = int(complexity)
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoHighCyclomaticComplexityRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		if hasGocycloIgnoreDirective(fn, file) {
			continue
		}

		c := cyclomaticComplexity(fn)
		if c > r.minComplexity {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Category:   lint.FailureCategoryComplexity,
				Failure:    fmt.Sprintf("function %s has cyclomatic complexity %d (> max enabled %d)", funcName(fn), c, r.minComplexity),
				Node:       fn,
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoHighCyclomaticComplexityRule) Name() string {
	return "noHighCyclomaticComplexity"
}

// Group returns the rule group.
func (*NoHighCyclomaticComplexityRule) Group() string {
	return "complexity"
}

// hasGocycloIgnoreDirective checks if a function declaration has a //gocyclo:ignore
// comment directive directly above it.
func hasGocycloIgnoreDirective(fn *ast.FuncDecl, file *lint.File) bool {
	if fn.Doc == nil {
		return false
	}
	for _, comment := range fn.Doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if text == "gocyclo:ignore" {
			return true
		}
	}
	return false
}

// cyclomaticComplexity calculates the cyclomatic complexity of a function.
// It starts with a base of 1 and increments for each if, for, case, &&, or ||.
func cyclomaticComplexity(fn *ast.FuncDecl) int {
	v := &cyclomaticComplexityVisitor{}
	ast.Walk(v, fn)
	return v.complexity
}

type cyclomaticComplexityVisitor struct {
	complexity int
}

// Visit implements the [ast.Visitor] interface.
func (v *cyclomaticComplexityVisitor) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.FuncDecl, *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt:
		v.complexity++
	case *ast.CaseClause:
		if n.List != nil {
			v.complexity++
		}
	case *ast.CommClause:
		if n.Comm != nil {
			v.complexity++
		}
	case *ast.BinaryExpr:
		if n.Op == token.LAND || n.Op == token.LOR {
			v.complexity++
		}
	}
	return v
}

// funcName returns the name representation of a function or method:
// "(Type).Name" for methods or simply "Name" for functions.
func funcName(fn *ast.FuncDecl) string {
	declarationHasReceiver := fn.Recv != nil && fn.Recv.NumFields() > 0
	if declarationHasReceiver {
		typ := fn.Recv.List[0].Type
		return fmt.Sprintf("(%s).%s", recvString(typ), fn.Name)
	}

	return fn.Name.Name
}

// recvString returns a string representation of recv of the
// form "T", "*T", or "BADRECV" (if not a proper receiver type).
func recvString(recv ast.Expr) string {
	switch t := recv.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + recvString(t.X)
	}
	return "BADRECV"
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}

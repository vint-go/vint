package no_excessive_results

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoExcessiveResultsRule detects functions with too many return values.
type NoExcessiveResultsRule struct {
	maxResults int
}

const defaultMaxResults = 5

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoExcessiveResultsRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.maxResults = defaultMaxResults
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		// try direct int64 argument
		results, ok := lint.ToInt64(arguments[0])
		if !ok {
			return fmt.Errorf(`invalid argument to the "noExcessiveResults" rule, expecting a k,v map or integer, got %T`, arguments[0])
		}
		r.maxResults = int(results)
		return nil
	}

	r.maxResults = defaultMaxResults
	for k, v := range argKV {
		if isRuleOption(k, "maxResults") {
			results, ok := lint.ToInt64(v)
			if !ok {
				return fmt.Errorf(`invalid configuration value for maxResults in "noExcessiveResults" rule; need integer but got %T`, v)
			}
			r.maxResults = int(results)
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoExcessiveResultsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	maxResults := r.maxResults
	if maxResults == 0 {
		maxResults = defaultMaxResults
	}

	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if funcDecl.Type.Results == nil {
			continue
		}

		numResults := funcDecl.Type.Results.NumFields()
		if numResults > maxResults {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Category:   lint.FailureCategoryComplexity,
				Failure:    fmt.Sprintf("function %s has too many results (%d > %d), consider grouping into a struct", funcName(funcDecl), numResults, maxResults),
				Node:       funcDecl,
			})
		}
	}

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoExcessiveResultsRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	maxResults := r.maxResults
	if maxResults == 0 {
		maxResults = defaultMaxResults
	}

	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return nil
	}

	if funcDecl.Type.Results == nil {
		return nil
	}

	numResults := funcDecl.Type.Results.NumFields()
	if numResults > maxResults {
		return []lint.Failure{
			{
				Confidence: 1,
				Category:   lint.FailureCategoryComplexity,
				Failure:    fmt.Sprintf("function %s has too many results (%d > %d), consider grouping into a struct", funcName(funcDecl), numResults, maxResults),
				Node:       funcDecl,
			},
		}
	}

	return nil
}

// Name returns the rule name.
func (*NoExcessiveResultsRule) Name() string {
	return "noExcessiveResults"
}

// Group returns the rule group.
func (*NoExcessiveResultsRule) Group() string {
	return "complexity"
}

// CacheTier returns the cache tier for this rule.
func (*NoExcessiveResultsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
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

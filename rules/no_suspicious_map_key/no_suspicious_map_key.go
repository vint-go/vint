package no_suspicious_map_key

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSuspiciousMapKeyRule detects suspicious map literal keys including
// whitespace anomalies in string keys and duplicate non-literal keys.
type NoSuspiciousMapKeyRule struct{}

// Apply applies the rule to given file.
func (r *NoSuspiciousMapKeyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}
	w := &lintSuspiciousMapKey{onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSuspiciousMapKeyRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}
	w := &lintSuspiciousMapKey{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoSuspiciousMapKeyRule) Name() string {
	return "noSuspiciousMapKey"
}

// Group returns the rule group.
func (*NoSuspiciousMapKeyRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoSuspiciousMapKeyRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintSuspiciousMapKey struct {
	onFailure func(lint.Failure)
}

func (w *lintSuspiciousMapKey) Visit(node ast.Node) ast.Visitor {
	compLit, ok := node.(*ast.CompositeLit)
	if !ok {
		return w
	}

	// Check if this is a map literal with string keys.
	mapType, ok := compLit.Type.(*ast.MapType)
	if !ok {
		return w
	}

	isStringKeyedMap := isStringType(mapType.Key)

	// Track keys for duplicate detection.
	seenKeys := map[string]bool{}

	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		// Check for whitespace anomalies in string keys.
		if isStringKeyedMap {
			if lit, ok := kv.Key.(*ast.BasicLit); ok {
				val := lit.Value
				// Strip the surrounding quotes.
				if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
					inner := val[1 : len(val)-1]
					if len(inner) > 0 {
						if strings.HasPrefix(inner, " ") && !strings.HasPrefix(inner, "  ") {
							w.onFailure(lint.Failure{
								Confidence: 1,
								Node:       kv.Key,
								Category:   lint.FailureCategoryLogic,
								Failure:    fmt.Sprintf("suspicious map key %s has leading space", val),
							})
						}
						if strings.HasSuffix(inner, " ") && !strings.HasSuffix(inner, "  ") {
							w.onFailure(lint.Failure{
								Confidence: 1,
								Node:       kv.Key,
								Category:   lint.FailureCategoryLogic,
								Failure:    fmt.Sprintf("suspicious map key %s has trailing space", val),
							})
						}
					}
				}
			}
		}

		// Check for duplicate keys (non-basic-literal expressions).
		if !isBasicLit(kv.Key) {
			keyStr := astutils.GoFmt(kv.Key)
			if seenKeys[keyStr] {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Node:       kv.Key,
					Category:   lint.FailureCategoryLogic,
					Failure:    fmt.Sprintf("suspicious duplicate key %s in map literal", keyStr),
				})
			} else {
				seenKeys[keyStr] = true
			}
		}
	}

	return w
}

// isStringType checks if the given expression refers to the built-in string type.
func isStringType(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "string"
}

// isBasicLit checks if the expression is a basic literal (string, int, float, etc.).
func isBasicLit(expr ast.Expr) bool {
	_, ok := expr.(*ast.BasicLit)
	return ok
}

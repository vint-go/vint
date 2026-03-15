package no_conflicting_http_mux_patterns

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoConflictingHttpMuxPatternsRule checks for conflicting HTTP mux patterns
// registered on the same ServeMux instance.
type NoConflictingHttpMuxPatternsRule struct{}

// Apply applies the rule to given file.
func (r *NoConflictingHttpMuxPatternsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintConflictingMux{
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

// Name returns the rule name.
func (*NoConflictingHttpMuxPatternsRule) Name() string {
	return "noConflictingHttpMuxPatterns"
}

// Group returns the rule group.
func (*NoConflictingHttpMuxPatternsRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoConflictingHttpMuxPatternsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type registeredPattern struct {
	pattern string
	node    ast.Node
}

type lintConflictingMux struct {
	onFailure func(lint.Failure)
}

func (w *lintConflictingMux) checkFunctionBody(body *ast.BlockStmt) {
	// Map from mux variable name to registered patterns
	muxPatterns := map[string][]registeredPattern{}

	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		muxName, pattern, ok := extractMuxPattern(call)
		if !ok {
			return true
		}

		// Check for conflicts with existing patterns on the same mux
		existing := muxPatterns[muxName]
		for _, ep := range existing {
			if patternsConflict(ep.pattern, pattern) {
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       call,
					Failure:    fmt.Sprintf("pattern %q conflicts with previously registered pattern %q", pattern, ep.pattern),
				})
				return true
			}
		}

		muxPatterns[muxName] = append(existing, registeredPattern{
			pattern: pattern,
			node:    call,
		})

		return true
	})
}

// extractMuxPattern extracts the mux variable name and pattern from a HandleFunc/Handle call.
// Returns (muxName, pattern, ok).
func extractMuxPattern(call *ast.CallExpr) (string, string, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", "", false
	}

	methodName := sel.Sel.Name
	if methodName != "HandleFunc" && methodName != "Handle" {
		return "", "", false
	}

	if len(call.Args) < 1 {
		return "", "", false
	}

	// Get pattern from first argument (must be a string literal)
	patternLit, ok := call.Args[0].(*ast.BasicLit)
	if !ok {
		return "", "", false
	}
	// Remove quotes from the string literal
	pattern := strings.Trim(patternLit.Value, `"`)

	// Get the mux variable name
	muxName := ""
	switch x := sel.X.(type) {
	case *ast.Ident:
		muxName = x.Name
	default:
		// For more complex expressions, use a generic key
		muxName = "_complex_expr_"
	}

	return muxName, pattern, true
}

// patternsConflict checks if two Go 1.22+ ServeMux patterns conflict.
// Two patterns conflict if they would match the same request.
func patternsConflict(p1, p2 string) bool {
	if p1 == p2 {
		return false // identical patterns are not conflicting in this context (duplicate, not conflict)
	}

	// Strip method prefix if present (e.g., "GET /items/{id}" -> "/items/{id}")
	method1, path1 := splitMethodAndPath(p1)
	method2, path2 := splitMethodAndPath(p2)

	// If both have methods and they differ, no conflict
	if method1 != "" && method2 != "" && method1 != method2 {
		return false
	}

	// Split paths into segments
	segs1 := splitPath(path1)
	segs2 := splitPath(path2)

	// Different number of segments means no conflict (unless one ends with {$} or wildcard...)
	if len(segs1) != len(segs2) {
		return false
	}

	// Check each segment pair
	for i := range segs1 {
		s1 := segs1[i]
		s2 := segs2[i]

		isWild1 := isWildcard(s1)
		isWild2 := isWildcard(s2)

		if isWild1 && isWild2 {
			// Two wildcards in the same position with different names conflict
			if s1 != s2 {
				return true
			}
			continue
		}

		if isWild1 || isWild2 {
			// One wildcard and one literal in the same position: no direct conflict
			// because the literal is more specific
			continue
		}

		// Both are literal segments
		if s1 != s2 {
			return false // different literal segments, no conflict
		}
	}

	return false
}

// splitMethodAndPath splits a pattern like "GET /items/{id}" into ("GET", "/items/{id}").
func splitMethodAndPath(pattern string) (string, string) {
	// Go 1.22 format: "METHOD /path" or just "/path"
	if idx := strings.Index(pattern, " /"); idx >= 0 {
		return pattern[:idx], pattern[idx+1:]
	}
	return "", pattern
}

// splitPath splits a URL path into segments.
func splitPath(path string) []string {
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}

// isWildcard returns true if the segment is a Go 1.22 wildcard like "{id}" or "{name...}".
func isWildcard(segment string) bool {
	return strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}")
}

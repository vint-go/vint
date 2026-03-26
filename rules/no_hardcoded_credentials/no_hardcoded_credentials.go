package no_hardcoded_credentials

import (
	"fmt"
	"go/ast"
	"go/token"
	"math"
	"regexp"
	"strconv"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoHardcodedCredentialsRule detects hardcoded credentials such as passwords,
// API keys, and tokens embedded directly in Go source code.
type NoHardcodedCredentialsRule struct {
	pattern          *regexp.Regexp
	entropyThreshold float64
}

// Configure implements lint.ConfigurableRule.
func (r *NoHardcodedCredentialsRule) Configure(arguments lint.Arguments) error {
	// Set defaults
	r.pattern = credentialNamePatterns
	r.entropyThreshold = defaultEntropyThreshold

	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noHardcodedCredentials" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch k {
		case "pattern":
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf(`invalid configuration value for pattern in "noHardcodedCredentials" rule; need string but got %T`, v)
			}
			re, err := regexp.Compile(s)
			if err != nil {
				return fmt.Errorf(`invalid regex pattern in "noHardcodedCredentials" rule: %w`, err)
			}
			r.pattern = re
		case "entropyThreshold":
			switch val := v.(type) {
			case float64:
				r.entropyThreshold = val
			case int64:
				r.entropyThreshold = float64(val)
			case string:
				f, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return fmt.Errorf(`invalid entropyThreshold in "noHardcodedCredentials" rule: %w`, err)
				}
				r.entropyThreshold = f
			default:
				return fmt.Errorf(`invalid configuration value for entropyThreshold in "noHardcodedCredentials" rule; need float64 but got %T`, v)
			}
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoHardcodedCredentialsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	pattern := r.pattern
	if pattern == nil {
		pattern = credentialNamePatterns
	}
	threshold := r.entropyThreshold
	if threshold == 0 {
		threshold = defaultEntropyThreshold
	}
	w := &lintHardcodedCredentials{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
		pattern:          pattern,
		entropyThreshold: threshold,
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoHardcodedCredentialsRule) Name() string {
	return "noHardcodedCredentials"
}

// Group returns the rule group.
func (*NoHardcodedCredentialsRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoHardcodedCredentialsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// credentialNamePatterns matches variable names that look like credentials.
var credentialNamePatterns = regexp.MustCompile(
	`(?i)(passwd|password|pwd|secret|token|api[_-]?key|apikey|access[_-]?key|auth[_-]?token|credentials?)`,
)

// Known secret format patterns.
var knownSecretPatterns = []*regexp.Regexp{
	// AWS Access Key ID (starts with AKIA, 20 chars)
	regexp.MustCompile(`^AKIA[0-9A-Z]{16}$`),
	// Slack token
	regexp.MustCompile(`^xox[bpras]-[0-9a-zA-Z-]+$`),
	// GitHub personal access token (classic or fine-grained)
	regexp.MustCompile(`^gh[ps]_[A-Za-z0-9_]{36,}$`),
	// Google API key
	regexp.MustCompile(`^AIza[0-9A-Za-z\-_]{35}$`),
}

// minEntropyLength is the minimum string length to consider for entropy analysis.
const minEntropyLength = 12

// defaultEntropyThreshold is the minimum Shannon entropy for a string to be flagged.
const defaultEntropyThreshold = 3.5

type lintHardcodedCredentials struct {
	onFailure        func(lint.Failure)
	pattern          *regexp.Regexp
	entropyThreshold float64
}

func (w *lintHardcodedCredentials) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.AssignStmt:
		w.checkAssignment(n)
	case *ast.GenDecl:
		w.checkGenDecl(n)
	case *ast.BinaryExpr:
		w.checkBinaryExpr(n)
	case *ast.CompositeLit:
		w.checkCompositeLit(n)
	}
	return w
}

// checkAssignment checks assignments like: password := "secret"
func (w *lintHardcodedCredentials) checkAssignment(stmt *ast.AssignStmt) {
	for i, lhs := range stmt.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok {
			continue
		}
		if !w.pattern.MatchString(ident.Name) {
			continue
		}
		if i >= len(stmt.Rhs) {
			continue
		}
		rhs := stmt.Rhs[i]
		if w.isHardcodedString(rhs) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       stmt,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "hardcoded credential: avoid embedding credentials in source code",
			})
		}
	}
}

// checkGenDecl checks const and var declarations: const apiKey = "..."
func (w *lintHardcodedCredentials) checkGenDecl(decl *ast.GenDecl) {
	if decl.Tok != token.CONST && decl.Tok != token.VAR {
		return
	}
	for _, spec := range decl.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for i, name := range vs.Names {
			if !w.pattern.MatchString(name.Name) {
				continue
			}
			if i >= len(vs.Values) {
				continue
			}
			if w.isHardcodedString(vs.Values[i]) {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Node:       vs,
					Category:   lint.FailureCategoryBadPractice,
					Failure:    "hardcoded credential: avoid embedding credentials in source code",
				})
			}
		}
	}
}

// checkBinaryExpr checks comparisons like: if userPassword == "admin123"
func (w *lintHardcodedCredentials) checkBinaryExpr(expr *ast.BinaryExpr) {
	if expr.Op != token.EQL && expr.Op != token.NEQ {
		return
	}

	// Check if left side is a credential-like name and right side is a string literal
	if ident, ok := expr.X.(*ast.Ident); ok {
		if w.pattern.MatchString(ident.Name) && w.isHardcodedString(expr.Y) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       expr,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "hardcoded credential: avoid comparing credentials against hardcoded values",
			})
			return
		}
	}

	// Check if right side is a credential-like name and left side is a string literal
	if ident, ok := expr.Y.(*ast.Ident); ok {
		if w.pattern.MatchString(ident.Name) && w.isHardcodedString(expr.X) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       expr,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "hardcoded credential: avoid comparing credentials against hardcoded values",
			})
		}
	}
}

// checkCompositeLit checks struct literals like: Config{Token: "ghp_..."}
func (w *lintHardcodedCredentials) checkCompositeLit(lit *ast.CompositeLit) {
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		ident, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		if !w.pattern.MatchString(ident.Name) {
			continue
		}
		if w.isHardcodedString(kv.Value) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       kv,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "hardcoded credential: avoid embedding credentials in struct literals",
			})
		}
	}
}

// isHardcodedString checks if an expression is a non-empty string literal.
func (w *lintHardcodedCredentials) isHardcodedString(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return false
	}
	val, err := strconv.Unquote(lit.Value)
	if err != nil {
		return false
	}
	if val == "" {
		return false
	}

	// Check if the value matches known secret patterns
	for _, pattern := range knownSecretPatterns {
		if pattern.MatchString(val) {
			return true
		}
	}

	// Perform entropy analysis for longer strings
	if len(val) >= minEntropyLength && shannonEntropy(val) >= w.entropyThreshold {
		return true
	}

	// If none of the above checks matched, the string doesn't look like a real credential
	return false
}

// shannonEntropy calculates the Shannon entropy of a string.
func shannonEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}

	freq := make(map[rune]float64)
	for _, c := range s {
		freq[c]++
	}

	length := float64(len([]rune(s)))
	entropy := 0.0
	for _, count := range freq {
		p := count / length
		if p > 0 {
			entropy -= p * math.Log2(p)
		}
	}

	return entropy
}

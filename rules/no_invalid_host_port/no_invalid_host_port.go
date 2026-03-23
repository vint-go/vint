package no_invalid_host_port

import (
	"go/ast"
	"go/token"
	"net"
	"strings"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoInvalidHostPortRule detects invalid host:port pairs passed to net/http
// functions like net.Listen, net.Dial, http.ListenAndServe, etc.
type NoInvalidHostPortRule struct{}

// Apply applies the rule to given file.
func (r *NoInvalidHostPortRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInvalidHostPort{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoInvalidHostPortRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInvalidHostPort{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoInvalidHostPortRule) Name() string {
	return "noInvalidHostPort"
}

// Group returns the rule group.
func (*NoInvalidHostPortRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInvalidHostPortRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// funcSpec describes which argument index contains the address for a given function.
type funcSpec struct {
	pkg     string
	name    string
	argIndex int
}

// netFuncs lists the functions whose address argument should be a valid host:port.
var netFuncs = []funcSpec{
	// net package
	{pkg: "net", name: "Dial", argIndex: 1},
	{pkg: "net", name: "DialTimeout", argIndex: 1},
	{pkg: "net", name: "Listen", argIndex: 1},
	{pkg: "net", name: "ListenPacket", argIndex: 1},
	{pkg: "net", name: "ResolveIPAddr", argIndex: 1},
	{pkg: "net", name: "ResolveUDPAddr", argIndex: 1},
	{pkg: "net", name: "ResolveTCPAddr", argIndex: 1},
	// http package
	{pkg: "http", name: "ListenAndServe", argIndex: 0},
	{pkg: "http", name: "ListenAndServeTLS", argIndex: 0},
}

type lintInvalidHostPort struct {
	onFailure func(lint.Failure)
}

func (w *lintInvalidHostPort) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	for _, spec := range netFuncs {
		if !astutils.IsPkgDotName(ce.Fun, spec.pkg, spec.name) {
			continue
		}

		if len(ce.Args) <= spec.argIndex {
			return w
		}

		arg := ce.Args[spec.argIndex]
		lit, ok := arg.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return w
		}

		// Extract the string value (remove surrounding quotes)
		val := lit.Value
		if len(val) >= 2 {
			if val[0] == '"' && val[len(val)-1] == '"' {
				val = val[1 : len(val)-1]
			} else if val[0] == '`' && val[len(val)-1] == '`' {
				val = val[1 : len(val)-1]
			}
		}

		// Empty string is valid (means all interfaces, port 0 or OS-assigned)
		if val == "" {
			return w
		}

		// Validate as host:port
		if !isValidHostPort(val) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ce,
				Category:   lint.FailureCategoryLogic,
				Failure:    "invalid host:port pair \"" + val + "\" passed to " + spec.pkg + "." + spec.name,
			})
		}

		return w
	}

	return w
}

// isValidHostPort checks whether addr is a valid host:port pair.
// A valid host:port must be parseable by net.SplitHostPort.
func isValidHostPort(addr string) bool {
	// net.SplitHostPort expects host:port format
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}

	// Port must not be empty (net.SplitHostPort allows empty port but we consider it suspicious)
	// Actually, empty port is valid for some use cases (e.g., ":0" means any port),
	// but a completely missing colon is what we want to catch.
	// net.SplitHostPort already rejects addresses without a colon, so if we get here
	// it means the format is correct.

	// Check that the port, if non-empty, doesn't have leading/trailing whitespace
	if port != strings.TrimSpace(port) {
		return false
	}

	return true
}

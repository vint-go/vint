package use_idiomatic_naming

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/vint-go/vint/internal/astutils"
	internalrule "github.com/vint-go/vint/internal/rule"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// knownNameExceptions is a set of known exceptions that should not be flagged.
var knownNameExceptions = map[string]bool{
	"LastInsertId": true, // must match database/sql
	"kWh":          true,
}

// UseIdiomaticNamingRule checks that identifiers follow Go naming conventions:
// MixedCaps or mixedCaps (not underscores), and acronyms should be all caps.
type UseIdiomaticNamingRule struct {
	extraInitialisms   []string // additional initialisms beyond defaults (e.g. GRPC, AMQP)
	excludeInitialisms []string // default initialisms to skip checking
}

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *UseIdiomaticNamingRule) Configure(arguments lint.Arguments) error {
	if len(arguments) == 0 {
		return nil
	}

	opts, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid argument to useIdiomaticNaming rule: expecting a map, got %T", arguments[0])
	}

	if v, ok := opts["extraInitialisms"]; ok {
		list, err := toStringSlice(v, "extraInitialisms")
		if err != nil {
			return err
		}
		r.extraInitialisms = list
	}

	if v, ok := opts["excludeInitialisms"]; ok {
		list, err := toStringSlice(v, "excludeInitialisms")
		if err != nil {
			return err
		}
		r.excludeInitialisms = list
	}

	return nil
}

func toStringSlice(v any, name string) ([]string, error) {
	switch val := v.(type) {
	case []any:
		result := make([]string, 0, len(val))
		for _, item := range val {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("invalid %s value: expecting string, got %T", name, item)
			}
			result = append(result, s)
		}
		return result, nil
	case []string:
		return val, nil
	default:
		return nil, fmt.Errorf("invalid %s: expecting a slice, got %T", name, v)
	}
}

// Apply applies the rule to the given file.
func (r *UseIdiomaticNamingRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintIdiomaticNaming{
		file: file,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
		extraInitialisms:   r.extraInitialisms,
		excludeInitialisms: r.excludeInitialisms,
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseIdiomaticNamingRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintIdiomaticNaming{
		file:               file,
		onFailure:          onFailure,
		extraInitialisms:   r.extraInitialisms,
		excludeInitialisms: r.excludeInitialisms,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseIdiomaticNamingRule) Name() string {
	return "useIdiomaticNaming"
}

// Group returns the rule group.
func (*UseIdiomaticNamingRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseIdiomaticNamingRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintIdiomaticNaming struct {
	file               *lint.File
	onFailure          func(lint.Failure)
	extraInitialisms   []string
	excludeInitialisms []string
}

func (w *lintIdiomaticNaming) checkFieldList(fl *ast.FieldList, thing string) {
	if fl == nil {
		return
	}
	for _, f := range fl.List {
		for _, id := range f.Names {
			w.check(id, thing)
		}
	}
}

func (w *lintIdiomaticNaming) check(id *ast.Ident, thing string) {
	if id.Name == "_" {
		return
	}
	if knownNameExceptions[id.Name] {
		return
	}

	should := internalrule.Name(id.Name, w.excludeInitialisms, w.extraInitialisms, false)
	if id.Name == should {
		return
	}

	if len(id.Name) > 2 && strings.Contains(id.Name[1:], "_") {
		w.onFailure(lint.Failure{
			Failure:    fmt.Sprintf("don't use underscores in Go names; %s %s should be %s", thing, id.Name, should),
			Confidence: 0.9,
			Node:       id,
			Category:   lint.FailureCategoryNaming,
		})
		return
	}
	w.onFailure(lint.Failure{
		Failure:    fmt.Sprintf("%s %s should be %s", thing, id.Name, should),
		Confidence: 0.8,
		Node:       id,
		Category:   lint.FailureCategoryNaming,
	})
}

func (w *lintIdiomaticNaming) Visit(n ast.Node) ast.Visitor {
	switch v := n.(type) {
	case *ast.AssignStmt:
		if v.Tok == token.ASSIGN {
			return w
		}
		for _, exp := range v.Lhs {
			if id, ok := exp.(*ast.Ident); ok {
				w.check(id, "var")
			}
		}
	case *ast.FuncDecl:
		funcName := v.Name.Name
		if w.file.IsTest() &&
			(strings.HasPrefix(funcName, "Example") ||
				strings.HasPrefix(funcName, "Test") ||
				strings.HasPrefix(funcName, "Benchmark") ||
				strings.HasPrefix(funcName, "Fuzz")) {
			return w
		}

		thing := "func"
		if v.Recv != nil {
			thing = "method"
		}

		// Exclude naming warnings for functions that are exported to C but
		// not exported in the Go API.
		if ast.IsExported(v.Name.Name) || !astutils.IsCgoExported(v) {
			w.check(v.Name, thing)
		}

		w.checkFieldList(v.Type.Params, thing+" parameter")
		w.checkFieldList(v.Type.Results, thing+" result")
	case *ast.GenDecl:
		if v.Tok == token.IMPORT {
			return w
		}

		thing := v.Tok.String()
		for _, spec := range v.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				w.check(s.Name, thing)
			case *ast.ValueSpec:
				for _, id := range s.Names {
					w.check(id, thing)
				}
			}
		}
	case *ast.InterfaceType:
		// Do not check interface method names.
		// They are often constrained by the method names of concrete types.
		for _, x := range v.Methods.List {
			ft, ok := x.Type.(*ast.FuncType)
			if !ok { // might be an embedded interface name
				continue
			}
			w.checkFieldList(ft.Params, "interface method parameter")
			w.checkFieldList(ft.Results, "interface method result")
		}
	case *ast.RangeStmt:
		if v.Tok == token.ASSIGN {
			return w
		}
		if id, ok := v.Key.(*ast.Ident); ok {
			w.check(id, "range var")
		}
		if id, ok := v.Value.(*ast.Ident); ok {
			w.check(id, "range var")
		}
	case *ast.StructType:
		for _, f := range v.Fields.List {
			for _, id := range f.Names {
				w.check(id, "struct field")
			}
		}
	}
	return w
}

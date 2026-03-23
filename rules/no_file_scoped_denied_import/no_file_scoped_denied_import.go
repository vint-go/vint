package no_file_scoped_denied_import

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/vint-go/vint/lint"
)

// NoFileScopedDeniedImportRule reports when a package import violates a
// file-scoped depguard rule, meaning the import is prohibited specifically
// in certain types of files based on glob patterns.
type NoFileScopedDeniedImportRule struct {
	ruleGroups map[string]*fileScopedDeniedImportGroup
}

type fileScopedDeniedImportGroup struct {
	name  string
	files []string
	deny  []fileScopedDeniedImportEntry
	allow []string // plain package patterns or "$gostd"
}

type fileScopedDeniedImportEntry struct {
	pkg  string
	desc string
}

// Configure validates and applies the rule configuration.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoFileScopedDeniedImportRule) Configure(arguments lint.Arguments) error {
	r.ruleGroups = map[string]*fileScopedDeniedImportGroup{}

	if len(arguments) == 0 {
		return nil
	}

	args, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid argument to the noFileScopedDeniedImport rule. Expecting a map, got %T", arguments[0])
	}

	// The config can be either a "rules" key containing the groups, or groups directly
	rulesMap := args
	if rVal, ok := args["rules"]; ok {
		rm, ok := rVal.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid 'rules' value in noFileScopedDeniedImport rule. Expecting a map, got %T", rVal)
		}
		rulesMap = rm
	}

	for groupName, groupCfg := range rulesMap {
		groupMap, ok := groupCfg.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid rule group %q configuration. Expecting a map, got %T", groupName, groupCfg)
		}

		group := &fileScopedDeniedImportGroup{
			name: groupName,
		}

		if filesVal, ok := groupMap["files"]; ok {
			if filesList, ok := filesVal.([]any); ok {
				for _, f := range filesList {
					if fStr, ok := f.(string); ok {
						group.files = append(group.files, fStr)
					}
				}
			}
		}

		if denyVal, ok := groupMap["deny"]; ok {
			denyList, ok := denyVal.([]any)
			if !ok {
				return fmt.Errorf("invalid deny list in rule group %q. Expecting a list, got %T", groupName, denyVal)
			}
			for _, d := range denyList {
				switch dv := d.(type) {
				case map[string]any:
					entry := fileScopedDeniedImportEntry{}
					if pkg, ok := dv["pkg"]; ok {
						if pkgStr, ok := pkg.(string); ok {
							entry.pkg = pkgStr
						}
					}
					if desc, ok := dv["desc"]; ok {
						if descStr, ok := desc.(string); ok {
							entry.desc = descStr
						}
					}
					group.deny = append(group.deny, entry)
				case string:
					group.deny = append(group.deny, fileScopedDeniedImportEntry{pkg: dv})
				}
			}
		}

		if allowVal, ok := groupMap["allow"]; ok {
			allowList, ok := allowVal.([]any)
			if !ok {
				return fmt.Errorf("invalid allow list in rule group %q. Expecting a list, got %T", groupName, allowVal)
			}
			for _, a := range allowList {
				switch av := a.(type) {
				case string:
					group.allow = append(group.allow, av)
				case map[string]any:
					if pkg, ok := av["pkg"]; ok {
						if pkgStr, ok := pkg.(string); ok {
							group.allow = append(group.allow, pkgStr)
						}
					}
				}
			}
		}

		r.ruleGroups[groupName] = group
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoFileScopedDeniedImportRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, importSpec := range file.AST.Imports {
		if importSpec.Path == nil {
			continue
		}
		importPath := strings.Trim(importSpec.Path.Value, `"`)

		for _, group := range r.ruleGroups {
			if !r.fileMatchesGroup(file, group) {
				continue
			}

			if entry, denied := r.isDenied(importPath, group); denied {
				msg := fmt.Sprintf("import '%s' is not allowed from list '%s'", importPath, group.name)
				if entry.desc != "" {
					msg += ": " + entry.desc
				}
				failures = append(failures, lint.Failure{
					Confidence: 1,
					Failure:    msg,
					Node:       importSpec,
					Category:   lint.FailureCategoryImports,
				})
			}
		}
	}

	return failures
}

// isDenied checks if an import is denied by the group's rules.
// If the group has a deny list, the import is denied if it matches any deny entry
// and is not covered by the allow list.
func (r *NoFileScopedDeniedImportRule) isDenied(importPath string, group *fileScopedDeniedImportGroup) (fileScopedDeniedImportEntry, bool) {
	for _, entry := range group.deny {
		if fileScopedMatchesPackage(importPath, entry.pkg) {
			// Check if the import is allowed by the allow list
			if r.isAllowed(importPath, group.allow) {
				return fileScopedDeniedImportEntry{}, false
			}
			return entry, true
		}
	}
	return fileScopedDeniedImportEntry{}, false
}

// isAllowed checks if the import path matches any entry in the allow list.
func (r *NoFileScopedDeniedImportRule) isAllowed(importPath string, allowList []string) bool {
	for _, pattern := range allowList {
		if pattern == "$gostd" {
			if isGoStdLib(importPath) {
				return true
			}
			continue
		}
		if fileScopedMatchesPackage(importPath, pattern) {
			return true
		}
	}
	return false
}

// fileScopedMatchesPackage checks if importPath matches the package pattern.
// Uses prefix matching: "foo/bar" matches "foo/bar" and "foo/bar/baz".
func fileScopedMatchesPackage(importPath, pattern string) bool {
	return importPath == pattern || strings.HasPrefix(importPath, pattern+"/")
}

// isGoStdLib checks if the import path looks like a Go standard library package.
// Standard library packages do not contain a dot in the first path segment.
func isGoStdLib(importPath string) bool {
	firstSlash := strings.Index(importPath, "/")
	firstSegment := importPath
	if firstSlash >= 0 {
		firstSegment = importPath[:firstSlash]
	}
	return !strings.Contains(firstSegment, ".")
}

// fileMatchesGroup checks if the file matches any file pattern in the group.
// If no file patterns are specified, the group applies to all files.
func (r *NoFileScopedDeniedImportRule) fileMatchesGroup(file *lint.File, group *fileScopedDeniedImportGroup) bool {
	if len(group.files) == 0 {
		return true
	}

	hasPositivePattern := false

	for _, pattern := range group.files {
		negate := false
		p := pattern
		if strings.HasPrefix(p, "!") {
			negate = true
			p = p[1:]
		}

		if !negate {
			hasPositivePattern = true
		}

		matched := false
		switch p {
		case "$all":
			matched = true
		case "$test":
			matched = file.IsTest()
		default:
			// Use glob matching on the filename
			m, err := filepath.Match(p, filepath.Base(file.Name))
			if err == nil && m {
				matched = true
			}
			// Also try matching against the full path
			if !matched {
				m, err = filepath.Match(p, file.Name)
				if err == nil && m {
					matched = true
				}
			}
		}

		if negate {
			if matched {
				return false
			}
		} else if matched {
			return true
		}
	}

	// If only negated patterns were specified and none matched (excluded) the file,
	// the file is considered to match the group.
	if !hasPositivePattern {
		return true
	}

	return false
}

// Name returns the rule name.
func (*NoFileScopedDeniedImportRule) Name() string {
	return "noFileScopedDeniedImport"
}

// Group returns the rule group.
func (*NoFileScopedDeniedImportRule) Group() string {
	return "correctness"
}

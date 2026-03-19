package use_package_naming

import (
	"errors"
	"fmt"
	"go/ast"
	"path/filepath"
	"regexp"
	"strings"

	gopackages "golang.org/x/tools/go/packages"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/internal/syncset"
	"github.com/strowk/vint/lint"
)

// defaultBadNames is the list of "bad" package names from https://go.dev/blog/package-names#bad-package-names.
var defaultBadNames = map[string]struct{}{
	"common":     {},
	"interface":  {},
	"interfaces": {},
	"misc":       {},
	"type":       {},
	"types":      {},
	"util":       {},
	"utils":      {},
}

// extraBadNames is the list of additional "bad" package names that are not recommended.
var extraBadNames = map[string]struct{}{
	"api":           {},
	"helpers":       {},
	"miscellaneous": {},
	"models":        {},
	"shared":        {},
	"utilities":     {},
}

// commonStdNames is the list of standard library package names that are commonly used in Go programs.
// This list is based on the most popular standard library packages according to importedby tab in pkg.go.dev.
// For example, "http" imported by 1,705,800 times https://pkg.go.dev/net/http?tab=importedby
var commonStdNames = map[string]string{
	"bytes":    "bytes",
	"bufio":    "bufio",
	"flag":     "flag",
	"context":  "context",
	"errors":   "errors",
	"filepath": "path/filepath",
	"fmt":      "fmt",
	"http":     "net/http",
	"io":       "io",
	"ioutil":   "io/ioutil",
	"json":     "encoding/json",
	"log":      "log",
	"math":     "math",
	"net":      "net",
	"os":       "os",
	"strconv":  "strconv",
	"reflect":  "reflect",
	"regexp":   "regexp",
	"runtime":  "runtime",
	"sort":     "sort",
	"strings":  "strings",
	"sync":     "sync",
	"time":     "time",
	"url":      "net/url",
}

// nonPublicPackageSegments are package path segments that indicate the std package is not public.
var nonPublicPackageSegments = map[string]struct{}{
	"internal": {},
	"vendor":   {},
}

// forbiddenTopLevelNames is the set of forbidden top level package names.
var forbiddenTopLevelNames = map[string]struct{}{
	"pkg": {},
}

// PackageNamingRule is a rule that checks package names.
type PackageNamingRule struct {
	skipConventionNameCheck  bool           // if true - skip checks for package name conventions (e.g., no underscores, no MixedCaps etc.)
	conventionNameCheckRegex *regexp.Regexp // the regex used to check package name conventions

	skipTopLevelCheck bool // if true - skip checks for top level package names (e.g., "pkg")

	skipDefaultBadNameCheck bool                // if true - skip checks for default bad package names (e.g., "util", "misc" etc.)
	checkExtraBadName       bool                // if true - enable check for extra bad package names (e.g., "helpers", "models" etc.)
	userDefinedBadNames     map[string]struct{} // set of user defined bad package names

	skipCollisionWithCommonStd bool // if true - skip checks for collisions with common Go standard library package names (e.g., "http", "json", "rand" etc.)

	checkCollisionWithAllStd bool // if true - enable checks for collisions with all Go standard library package names (including "version", "metrics" etc.)
	// allStdNames holds name -> path of standard library packages excluding internal and vendor.
	// Populated only if checkCollisionWithAllStd is true. `net/http` stored as `http`, `math/rand/v2` as `rand` etc.
	allStdNames map[string]string

	// alreadyCheckedNames is keyed by fileDir (package directory path) to track which package directories
	// have already been checked and avoid duplicate checks across files in the same package.
	alreadyCheckedNames *syncset.Set
}

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *PackageNamingRule) Configure(arguments lint.Arguments) error {
	r.alreadyCheckedNames = syncset.New()

	if len(arguments) == 0 {
		return nil
	}

	if len(arguments) > 1 {
		return fmt.Errorf("invalid arguments to the usePackageNaming rule: expected at most 1 argument, but got %d", len(arguments))
	}

	args, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid argument to the usePackageNaming rule: expecting a k,v map, but got %T", arguments[0])
	}

	for k, v := range args {
		switch {
		case isRuleOption(k, "skipConventionNameCheck"):
			r.skipConventionNameCheck, ok = v.(bool)
			if !ok {
				return fmt.Errorf("invalid argument to the usePackageNaming rule: expecting skipConventionNameCheck to be a boolean, but got %T", v)
			}
		case isRuleOption(k, "conventionNameCheckRegex"):
			regexStr, ok := v.(string)
			if !ok {
				return fmt.Errorf("invalid argument to the usePackageNaming rule: expecting conventionNameCheckRegex to be a string, but got %T", v)
			}
			if regexStr == "" {
				return errors.New("invalid argument to the usePackageNaming rule: conventionNameCheckRegex cannot be an empty string")
			}
			regex, err := regexp.Compile(regexStr)
			if err != nil {
				return fmt.Errorf("invalid argument to the usePackageNaming rule: invalid regex for conventionNameCheckRegex: %w", err)
			}
			r.conventionNameCheckRegex = regex
		case isRuleOption(k, "skipTopLevelCheck"):
			r.skipTopLevelCheck, ok = v.(bool)
			if !ok {
				return fmt.Errorf("invalid argument to the usePackageNaming rule: expecting skipTopLevelCheck to be a boolean, but got %T", v)
			}
		case isRuleOption(k, "skipDefaultBadNameCheck"):
			r.skipDefaultBadNameCheck, ok = v.(bool)
			if !ok {
				return fmt.Errorf("invalid argument to the usePackageNaming rule: expecting skipDefaultBadNameCheck to be a boolean, but got %T", v)
			}
		case isRuleOption(k, "checkExtraBadName"):
			r.checkExtraBadName, ok = v.(bool)
			if !ok {
				return fmt.Errorf("invalid argument to the usePackageNaming rule: expecting checkExtraBadName to be a boolean, but got %T", v)
			}
		case isRuleOption(k, "userDefinedBadNames"):
			userDefinedBadNames, ok := v.([]any)
			if !ok {
				return fmt.Errorf("invalid argument to the usePackageNaming rule: expecting userDefinedBadNames of type slice of strings, but got %T", v)
			}
			for i, name := range userDefinedBadNames {
				if r.userDefinedBadNames == nil {
					r.userDefinedBadNames = map[string]struct{}{}
				}
				n, ok := name.(string)
				if !ok {
					return fmt.Errorf("invalid argument to the usePackageNaming rule: expecting element %d of userDefinedBadNames to be a string, but got %v(%T)", i, name, name)
				}
				if n == "" {
					return fmt.Errorf("invalid argument to the usePackageNaming rule: userDefinedBadNames cannot contain empty string (index %d)", i)
				}
				r.userDefinedBadNames[strings.ToLower(n)] = struct{}{}
			}
		case isRuleOption(k, "skipCollisionWithCommonStd"):
			r.skipCollisionWithCommonStd, ok = v.(bool)
			if !ok {
				return fmt.Errorf("invalid argument to the usePackageNaming rule: expecting skipCollisionWithCommonStd to be a boolean, but got %T", v)
			}
		case isRuleOption(k, "checkCollisionWithAllStd"):
			r.checkCollisionWithAllStd, ok = v.(bool)
			if !ok {
				return fmt.Errorf("invalid argument to the usePackageNaming rule: expecting checkCollisionWithAllStd to be a boolean, but got %T", v)
			}
		}
	}

	if r.skipConventionNameCheck && r.conventionNameCheckRegex != nil {
		return errors.New("invalid configuration for usePackageNaming rule: skipConventionNameCheck and conventionNameCheckRegex cannot be both set")
	}

	if r.skipCollisionWithCommonStd && r.checkCollisionWithAllStd {
		return errors.New("invalid configuration for usePackageNaming rule: skipCollisionWithCommonStd and checkCollisionWithAllStd cannot be both set")
	}

	if r.checkCollisionWithAllStd && r.allStdNames == nil {
		pkgs, err := gopackages.Load(nil, "std")
		if err != nil {
			return fmt.Errorf("load std packages: %w", err)
		}

		r.allStdNames = map[string]string{}
		for _, pkg := range pkgs {
			if isNonPublicPackage(pkg.PkgPath) {
				continue
			}
			if existingPath, ok := r.allStdNames[pkg.Name]; !ok || pkg.PkgPath < existingPath {
				r.allStdNames[pkg.Name] = pkg.PkgPath
			}
		}
	}

	return nil
}

// isNonPublicPackage reports whether the path represents an internal or vendor directory.
func isNonPublicPackage(path string) bool {
	for p := range strings.SplitSeq(path, "/") {
		if _, ok := nonPublicPackageSegments[p]; ok {
			return true
		}
	}
	return false
}

// Apply applies the rule to given file.
func (r *PackageNamingRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	fileDir := filepath.Dir(file.Name)

	if !r.alreadyCheckedNames.AddIfAbsent(fileDir) {
		return failures
	}

	node := file.AST.Name
	pkgName := node.Name
	pkgNameWithoutTestSuffix := strings.TrimSuffix(pkgName, "_test")

	if r.conventionNameCheckRegex != nil {
		if !r.conventionNameCheckRegex.MatchString(pkgNameWithoutTestSuffix) {
			onFailure(r.pkgNameFailure(node, "package name %q doesn't match the convention defined by conventionNameCheckRegex", pkgName))
			return failures
		}
	} else if !r.skipConventionNameCheck {
		// Package names need slightly different handling than other names.
		if strings.Contains(pkgNameWithoutTestSuffix, "_") {
			onFailure(r.pkgNameFailure(node, "don't use package name %q that contains an underscore", pkgName))
			return failures
		}
		if hasUpperCaseLetter(pkgNameWithoutTestSuffix) {
			onFailure(r.pkgNameFailure(node, "don't use package name %q that contains MixedCaps", pkgName))
			return failures
		}
	}

	pkgNameLower := strings.ToLower(pkgName)
	if !r.skipTopLevelCheck {
		if _, ok := forbiddenTopLevelNames[pkgNameLower]; ok && filepath.Base(fileDir) != pkgName {
			onFailure(r.pkgNameFailure(node, "don't use %q as a root level package name", pkgName))
			return failures
		}
	}

	if !r.skipDefaultBadNameCheck {
		if _, ok := defaultBadNames[pkgNameLower]; ok {
			onFailure(r.pkgNameFailure(node, "don't use %q because it is a bad package name according to https://go.dev/blog/package-names#bad-package-names", pkgName))
			return failures
		}
	}

	if r.checkExtraBadName {
		if _, ok := extraBadNames[pkgNameLower]; ok {
			onFailure(r.pkgNameFailure(node, "don't use %q because it is a bad package name (extra)", pkgName))
			return failures
		}
	}

	if r.userDefinedBadNames != nil {
		if _, ok := r.userDefinedBadNames[pkgNameLower]; ok {
			onFailure(r.pkgNameFailure(node, "don't use %q because it is a bad package name (user-defined)", pkgName))
			return failures
		}
	}

	if r.checkCollisionWithAllStd {
		// all std names are also common std names, so no need to check separately
		if std, ok := r.allStdNames[pkgNameLower]; ok {
			onFailure(r.pkgNameFailure(node, "don't use %q because it conflicts with Go standard library package %q", pkgName, std))
		}
	} else if !r.skipCollisionWithCommonStd {
		if std, ok := commonStdNames[pkgNameLower]; ok {
			onFailure(r.pkgNameFailure(node, "don't use %q because it conflicts with common Go standard library package %q", pkgName, std))
		}
	}

	return failures
}

// Name returns the rule name.
func (*PackageNamingRule) Name() string {
	return "usePackageNaming"
}

// Group returns the rule group.
func (*PackageNamingRule) Group() string {
	return "style"
}

func (*PackageNamingRule) pkgNameFailure(node ast.Node, msg string, args ...any) lint.Failure {
	return lint.Failure{
		Failure:    fmt.Sprintf(msg, args...),
		Confidence: 1,
		Node:       node,
		Category:   lint.FailureCategoryNaming,
	}
}

// CacheTier returns the cache tier for this rule.
func (*PackageNamingRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// hasUpperCaseLetter checks if a string contains at least one upper case letter.
func hasUpperCaseLetter(s string) bool {
	for _, r := range s {
		if isUpper(r) {
			return true
		}
	}
	return false
}

// isUpper checks if a rune is an uppercase letter.
func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}

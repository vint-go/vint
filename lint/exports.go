package lint

import (
	"bufio"
	"bytes"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	goversion "github.com/hashicorp/go-version"
	"github.com/strowk/vint/internal/rulecache"
	"golang.org/x/mod/modfile"
)

// SharedImporter is an exported alias for the shared importer type.
// It caches resolved packages across a lint run so each import path
// is resolved at most once.
type SharedImporter = sharedImporter

// NewSharedImporter creates a new shared importer suitable for reuse
// across multiple packages in a single lint run.
func NewSharedImporter() *SharedImporter {
	return newSharedImporter()
}

// NewPackage creates a new Package with the given FileSet, Go version, and
// shared importer. The package starts with no files — use AddFile to populate.
func NewPackage(fset *token.FileSet, goVersion *goversion.Version, imp *SharedImporter) *Package {
	return &Package{
		fset:      fset,
		importer:  imp,
		files:     map[string]*File{},
		goVersion: goVersion,
	}
}

// AddFile parses the given source and adds it as a file in this package.
// Returns the created File or an error if parsing fails.
func (p *Package) AddFile(name string, content []byte) (*File, error) {
	file, err := NewFile(name, content, p)
	if err != nil {
		return nil, err
	}
	p.mu.Lock()
	p.files[name] = file
	p.mu.Unlock()
	return file, nil
}

// FileSet returns the package's token.FileSet.
func (p *Package) FileSet() *token.FileSet {
	return p.fset
}

// ExportedScanSortable exposes scanSortable for external orchestrators.
func (p *Package) ExportedScanSortable() {
	p.scanSortable()
}

// DetectGoMod finds the nearest go.mod file starting from dir and walking up.
// Returns the module root directory and the Go version declared in go.mod.
func DetectGoMod(dir string) (rootDir string, ver *goversion.Version, err error) {
	modFileName, err := retrieveModFile(dir)
	if err != nil {
		return "", nil, fmt.Errorf("%q doesn't seem to be part of a Go module", dir)
	}

	mod, err := os.ReadFile(modFileName) //nolint:gosec
	if err != nil {
		return "", nil, fmt.Errorf("failed to read %q, got %w", modFileName, err)
	}

	modAst, err := modfile.ParseLax(modFileName, mod, nil)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse %q, got %w", modFileName, err)
	}

	if modAst.Go == nil {
		return "", nil, fmt.Errorf("%q does not specify a Go version", modFileName)
	}

	ver, err = goversion.NewVersion(modAst.Go.Version)
	return filepath.Dir(modFileName), ver, err
}

// DefaultGoVersion is the fallback Go version used when no go.mod is found.
var DefaultGoVersion = goversion.Must(goversion.NewVersion("1.0"))

// IsGenerated reports whether the source file is generated code
// according to the rules from https://go.dev/s/generatedcode.
func IsGenerated(src []byte) bool {
	sc := bufio.NewScanner(bytes.NewReader(src))
	for sc.Scan() {
		b := sc.Bytes()
		if bytes.HasPrefix(b, generatedPrefix) && bytes.HasSuffix(b, generatedSuffix) && len(b) >= len(generatedPrefix)+len(generatedSuffix) {
			return true
		}
	}
	return false
}

// AddInvalidFileFailure sends a failure for a file that could not be parsed.
func AddInvalidFileFailure(filename, errStr string, failures chan<- Failure) {
	position := getPositionInvalidFile(filename, errStr)
	failures <- Failure{
		Confidence: 1,
		Failure:    fmt.Sprintf("invalid file %s: %v", filename, errStr),
		Category:   failureCategoryValidity,
		Position:   position,
	}
}

// DisabledIntervalsForFile computes disabled intervals for a file from
// revive directives in comments. This is a standalone function that can be
// called from external packages.
func DisabledIntervalsForFile(f *File, rules []Rule, mustSpecifyDisableReason bool, failures chan<- Failure) map[string][]DisabledInterval {
	enabledDisabledRulesMap := map[string][]enableDisableConfig{}

	getEnabledDisabledIntervals := func() map[string][]DisabledInterval {
		result := map[string][]DisabledInterval{}
		for ruleName, disabledArr := range enabledDisabledRulesMap {
			ruleResult := []DisabledInterval{}
			for i := range disabledArr {
				interval := DisabledInterval{
					RuleName: ruleName,
					From: token.Position{
						Filename: f.Name,
						Line:     disabledArr[i].position,
					},
					To: token.Position{
						Filename: f.Name,
						Line:     1<<31 - 1, // math.MaxInt32
					},
				}
				if i%2 == 0 {
					ruleResult = append(ruleResult, interval)
				} else {
					ruleResult[len(ruleResult)-1].To.Line = disabledArr[i].position
				}
			}
			result[ruleName] = ruleResult
		}
		return result
	}

	handleConfig := func(isEnabled bool, line int, name string) {
		existing := enabledDisabledRulesMap[name]
		if (len(existing) > 1 && existing[len(existing)-1].enabled == isEnabled) ||
			(len(existing) == 0 && isEnabled) {
			return
		}
		enabledDisabledRulesMap[name] = append(existing, enableDisableConfig{
			enabled:  isEnabled,
			position: line,
		})
	}

	handleRules := func(modifier string, isEnabled bool, line int, ruleNames []string) {
		for _, name := range ruleNames {
			switch modifier {
			case "line":
				handleConfig(isEnabled, line, name)
				handleConfig(!isEnabled, line, name)
			case "next-line":
				handleConfig(isEnabled, line+1, name)
				handleConfig(!isEnabled, line+1, name)
			default:
				handleConfig(isEnabled, line, name)
			}
		}
	}

	for _, c := range f.AST.Comments {
		line := f.ToPosition(c.End()).Line
		for _, comment := range c.List {
			match := directiveRegexp.FindStringSubmatch(comment.Text)
			if len(match) == 0 {
				continue
			}
			ruleNames := []string{}
			for name := range strings.SplitSeq(match[rulesPos], ",") {
				name = strings.Trim(name, "\n")
				if name != "" {
					ruleNames = append(ruleNames, name)
				}
			}

			mustCheckDisablingReason := mustSpecifyDisableReason && match[directivePos] == "disable"
			if mustCheckDisablingReason && strings.Trim(match[reasonPos], " ") == "" {
				failures <- Failure{
					Confidence: 1,
					RuleName:   directiveSpecifyDisableReason,
					Failure:    "reason of lint disabling not found",
					Position:   ToFailurePosition(comment.Pos(), comment.End(), f),
					Node:       comment,
				}
				continue
			}

			if len(ruleNames) == 0 {
				for _, rule := range rules {
					ruleNames = append(ruleNames, FullRuleName(rule))
				}
			}

			handleRules(match[modifierPos], match[directivePos] == "enable", line, ruleNames)
		}
	}

	return getEnabledDisabledIntervals()
}

// ComputeConfigHash produces a single hash covering all rule configurations
// and global settings that affect lint output.
func ComputeConfigHash(ruleSet []Rule, config Config) [32]byte {
	return computeConfigHash(ruleSet, config)
}

// ToCachedFailure converts a Failure to a CachedFailure for cache storage.
func ToCachedFailure(f Failure) rulecache.CachedFailure {
	return toCachedFailure(f)
}

// FromCachedFailure converts a CachedFailure back to a Failure.
func FromCachedFailure(cf rulecache.CachedFailure) Failure {
	return fromCachedFailure(cf)
}

// CollectSiblingContents returns the content of all files in the same package as f.
func CollectSiblingContents(f *File) map[string][]byte {
	return f.collectSiblingContents()
}

// FilterFailuresByDisabledIntervals filters out failures that fall within
// disabled intervals. Standalone function for external packages.
func FilterFailuresByDisabledIntervals(failures []Failure, disabledIntervals map[string][]DisabledInterval) []Failure {
	result := make([]Failure, 0, len(failures))
	for _, failure := range failures {
		fStart := failure.Position.Start.Line
		fEnd := failure.Position.End.Line
		intervals, ok := disabledIntervals[failure.RuleName]
		if !ok {
			result = append(result, failure)
			continue
		}

		include := true
		for _, interval := range intervals {
			intStart := interval.From.Line
			intEnd := interval.To.Line
			if (fStart >= intStart && fStart <= intEnd) ||
				(fEnd >= intStart && fEnd <= intEnd) {
				include = false
				break
			}
		}
		if include {
			result = append(result, failure)
		}
	}
	return result
}

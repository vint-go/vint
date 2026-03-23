package no_duplicate_code

import (
	"fmt"
	"go/token"
	"strings"
	"sync"

	"github.com/vint-go/vint/lint"
)

const defaultThreshold = 150

// Compile-time interface checks.
var (
	_ lint.Rule             = (*NoDuplicateCodeRule)(nil)
	_ lint.AggregatingRule  = (*NoDuplicateCodeRule)(nil)
	_ lint.ConfigurableRule = (*NoDuplicateCodeRule)(nil)
	_ lint.Grouped          = (*NoDuplicateCodeRule)(nil)
)

// NoDuplicateCodeRule detects duplicate code across files using
// a suffix tree algorithm with GC-friendly flat data structures.
type NoDuplicateCodeRule struct {
	threshold int
	mu        sync.Mutex
	global    tokenData  // concatenated token data from all files
	spans     []fileSpan // per-file ranges into global
}

// Name returns the rule name.
func (*NoDuplicateCodeRule) Name() string { return "noDuplicateCode" }

// Group returns the rule group.
func (*NoDuplicateCodeRule) Group() string { return "complexity" }

// Configure validates and applies rule configuration.
func (r *NoDuplicateCodeRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.threshold = defaultThreshold
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		threshold, ok := lint.ToInt64(arguments[0])
		if !ok {
			return fmt.Errorf(`invalid argument to the "noDuplicateCode" rule, expecting a k,v map or integer, got %T`, arguments[0])
		}
		r.threshold = int(threshold)
		return nil
	}

	for k, v := range argKV {
		if normalizeOption(k) == normalizeOption("threshold") {
			threshold, ok := lint.ToInt64(v)
			if !ok || threshold < 0 {
				return fmt.Errorf(`invalid configuration value for threshold in "noDuplicateCode" rule; need positive integer but got %T`, v)
			}
			r.threshold = int(threshold)
		}
	}

	if r.threshold == 0 {
		r.threshold = defaultThreshold
	}
	return nil
}

// Apply returns nil — this rule produces results via Collect/Finalize.
func (r *NoDuplicateCodeRule) Apply(_ *lint.File, _ lint.Arguments) []lint.Failure {
	return nil
}

// Collect serializes a file's already-parsed AST into flat int32 arrays
// and appends to the global token sequence. No re-parsing from disk,
// no *syntax.Node heap allocations. Safe for concurrent calls.
func (r *NoDuplicateCodeRule) Collect(file *lint.File, _ lint.Arguments) {
	td := serializeAST(file.AST, file.Pkg.FileSet())

	r.mu.Lock()
	start := r.global.len()
	r.global.types = append(r.global.types, td.types...)
	r.global.pos = append(r.global.pos, td.pos...)
	r.global.end = append(r.global.end, td.end...)
	r.global.owns = append(r.global.owns, td.owns...)
	r.spans = append(r.spans, fileSpan{
		start: start,
		end:   r.global.len(),
		file:  file.Name,
	})
	r.mu.Unlock()
}

// Finalize builds a suffix tree from all collected token data,
// finds duplicate code fragments, and returns failures.
func (r *NoDuplicateCodeRule) Finalize() []lint.Failure {
	if r.threshold == 0 {
		r.threshold = defaultThreshold
	}

	if len(r.spans) < 2 {
		return nil
	}

	// Build suffix tree on the flat int32 type sequence.
	tree := newSTree(len(r.global.types) + 1)
	tree.update(r.global.types)

	// Add sentinel to finalize the tree.
	tree.update([]int32{nodeSentinel})

	// Find all duplicate substrings above threshold.
	matches := findDuplicates(tree, r.threshold)

	// Deduplicate: track reported (file, pos) pairs.
	type fragKey struct {
		file string
		pos  int32
	}
	reported := map[[2]fragKey]bool{}
	var failures []lint.Failure

	for _, m := range matches {
		sm := findSyntaxUnits(&r.global, m, r.threshold)
		if len(sm.frags) < 2 {
			continue
		}

		// Report each pair of fragment occurrences.
		for i := 0; i < len(sm.frags); i++ {
			for j := i + 1; j < len(sm.frags); j++ {
				fragI := sm.frags[i]
				fragJ := sm.frags[j]
				if len(fragI) == 0 || len(fragJ) == 0 {
					continue
				}

				fileI := r.fileForIndex(fragI[0].startIdx)
				fileJ := r.fileForIndex(fragJ[0].startIdx)

				posI := r.global.pos[fragI[0].startIdx]
				posJ := r.global.pos[fragJ[0].startIdx]

				keyI := fragKey{fileI, posI}
				keyJ := fragKey{fileJ, posJ}

				// Normalize order for dedup key.
				pairKey := [2]fragKey{keyI, keyJ}
				if fileI > fileJ || (fileI == fileJ && posI > posJ) {
					pairKey = [2]fragKey{keyJ, keyI}
				}

				if reported[pairKey] {
					continue
				}
				reported[pairKey] = true

				lastFragJ := fragJ[len(fragJ)-1]

				failures = append(failures, lint.Failure{
					Confidence: 1,
					Category:   lint.FailureCategoryComplexity,
					Failure: fmt.Sprintf(
						"duplicate code detected: %s has code duplicated from %s (>= %d tokens of identical structure)",
						fileJ, fileI, r.threshold,
					),
					Position: lint.FailurePosition{
						Start: token.Position{
							Filename: fileJ,
							Offset:   int(posJ),
						},
						End: token.Position{
							Filename: fileJ,
							Offset:   int(r.global.end[lastFragJ.startIdx]),
						},
					},
				})
			}
		}
	}

	return failures
}

// fileForIndex returns the filename for a given global token index.
func (r *NoDuplicateCodeRule) fileForIndex(idx int) string {
	for _, span := range r.spans {
		if idx >= span.start && idx < span.end {
			return span.file
		}
	}
	return "<unknown>"
}

func normalizeOption(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "-", ""))
}

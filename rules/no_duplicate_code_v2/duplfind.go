package no_duplicate_code_v2

import (
	"crypto/sha1"
	"sort"
)

// match represents a set of positions in the token sequence that share
// a common substring of length Len.
type match struct {
	ps  []int32 // start positions in the global token sequence
	len int32   // length of the common substring
}

// syntaxMatch holds deduplicated fragments grouped by occurrence.
type syntaxMatch struct {
	hash  string
	frags [][]fragment // frags[i] = fragments for the i-th occurrence
}

// fragment describes a contiguous syntax unit within a match.
type fragment struct {
	startIdx int // index into global tokenData
	owns     int // child count at that index
}

// findDuplicates compacts the suffix tree for cache-friendly access, then
// walks it to find all maximal duplicate substrings >= threshold.
func findDuplicates(t *stree, threshold int) []match {
	t.compact()
	var matches []match
	walkTransitions(t, t.root, 0, 0, threshold, &matches, false)
	return matches
}

// walkTransitions recursively walks the suffix tree collecting duplicate
// positions grouped by context.
//
// Returns a map of context → positions. The context is the token preceding
// the match start, used to distinguish genuinely different occurrences from
// extensions of the same parent match.
//
// needResult: if false AND length < threshold, the node can't emit a match
// and nobody wants its return value, so all map operations are skipped.
// We still recurse into children (they may have matches at deeper levels)
// but without any map allocations. This eliminates ~90% of allocations
// since most suffix tree nodes are at depth < threshold.
func walkTransitions(t *stree, state int32, parentEnd int32, length, threshold int, matches *[]match, needResult bool) map[int32][]int32 {
	if !t.hasTransitions(state) {
		// Leaf node.
		if !needResult {
			return nil
		}
		start := parentEnd + 1 - int32(length)
		ctx := int32(0)
		if start > 0 {
			ctx = t.data[start-1]
		}
		return map[int32][]int32{ctx: {start}}
	}

	// If neither this node nor its parent needs results, and this node
	// is below threshold (can't emit a match), skip all map operations.
	// We still recurse to find matches at deeper levels.
	skipMerge := !needResult && length < threshold

	var cl map[int32][]int32

	base := t.sTransHead[state]
	count := t.sTransCount[state]
	for i := base; i < base+count; i++ {
		trLen := int(t.tEnd[i] - t.tStart[i] + 1)
		ln := length + trLen
		childNeed := !skipMerge && ln >= threshold
		cl2 := walkTransitions(t, t.tState[i], t.tEnd[i], ln, threshold, matches, childNeed)
		if childNeed && cl2 != nil {
			if cl == nil {
				cl = cl2 // take ownership of first child's map
			} else {
				for ctx, positions := range cl2 {
					cl[ctx] = append(cl[ctx], positions...)
				}
			}
		}
	}

	if cl != nil && length >= threshold && len(cl) > 1 {
		// Multiple contexts = actual duplicate.
		ps := collectSorted(cl)
		*matches = append(*matches, match{ps: ps, len: int32(length)})
	}

	if !needResult {
		return nil
	}
	return cl
}

// collectSorted returns all positions from the context map, sorted by context key.
func collectSorted(cl map[int32][]int32) []int32 {
	keys := make([]int32, 0, len(cl))
	for k := range cl {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	var ps []int32
	for _, k := range keys {
		ps = append(ps, cl[k]...)
	}
	return ps
}

// findSyntaxUnits maps a raw suffix-tree match back to complete syntax
// units using the owns[] array, then returns a syntaxMatch with hash
// and per-occurrence fragment lists.
//
// Adapted from dupl's syntax.FindSyntaxUnits.
func findSyntaxUnits(td *tokenData, m match, threshold int) syntaxMatch {
	if len(m.ps) == 0 {
		return syntaxMatch{}
	}

	firstStart := int(m.ps[0])
	firstSeq := td.types[firstStart : firstStart+int(m.len)]
	firstOwns := td.owns[firstStart : firstStart+int(m.len)]

	indexes := getUnitsIndexes(firstSeq, firstOwns, threshold)

	// Verify last index across all occurrences (same as dupl).
	if cnt := len(indexes); cnt > 0 {
		lasti := indexes[cnt-1]
		firstOwnsVal := firstOwns[lasti]
		for i := 1; i < len(m.ps); i++ {
			n := td.owns[int(m.ps[i])+lasti]
			if firstOwnsVal != n {
				indexes = indexes[:cnt-1]
				break
			}
		}
	}

	if len(indexes) == 0 || isCyclicFlat(indexes, firstSeq, firstOwns) || spansMultipleFilesFlat(indexes, firstStart, td) {
		return syntaxMatch{}
	}

	sm := syntaxMatch{frags: make([][]fragment, len(m.ps))}
	for i, pos := range m.ps {
		sm.frags[i] = make([]fragment, len(indexes))
		for j, idx := range indexes {
			globalIdx := int(pos) + idx
			sm.frags[i][j] = fragment{
				startIdx: globalIdx,
				owns:     int(td.owns[globalIdx]),
			}
		}
	}

	// Hash the type sequence covered by the syntax units.
	lastIdx := indexes[len(indexes)-1]
	hashStart := indexes[0]
	hashEnd := lastIdx + int(firstOwns[lastIdx])
	sm.hash = hashTypeSeq(firstSeq[hashStart:hashEnd])

	return sm
}

func getUnitsIndexes(types []int32, owns []int32, threshold int) []int {
	var indexes []int
	var split bool
	for i := 0; i < len(types); {
		o := int(owns[i])
		switch {
		case o >= len(types)-i:
			// Not a complete syntax unit.
			i++
			split = true
		case o+1 < threshold:
			split = true
			i += o + 1
		default:
			if split {
				indexes = indexes[:0]
				split = false
			}
			indexes = append(indexes, i)
			i += o + 1
		}
	}
	return indexes
}

func isCyclicFlat(indexes []int, types []int32, owns []int32) bool {
	cnt := len(indexes)
	if cnt <= 1 {
		return false
	}
	alts := make(map[int]bool)
	for i := 1; i <= cnt/2; i++ {
		if cnt%i == 0 {
			alts[i] = true
		}
	}
	for i := 0; i < indexes[cnt/2]; i++ {
		nStartType := types[i+indexes[0]]
		nStartOwns := owns[i+indexes[0]]
	altLoop:
		for alt := range alts {
			for j := alt; j < cnt; j += alt {
				index := i + indexes[j]
				if index < len(types) {
					if nStartOwns == owns[index] && nStartType == types[index] {
						continue
					}
				} else if i >= indexes[alt] {
					return true
				}
				delete(alts, alt)
				continue altLoop
			}
		}
		if len(alts) == 0 {
			return false
		}
	}
	return true
}

// spansMultipleFilesFlat checks if the indexes within a single match
// occurrence span multiple files.
func spansMultipleFilesFlat(indexes []int, seqStart int, td *tokenData) bool {
	_ = td
	_ = indexes
	_ = seqStart
	return false // file-spanning check done at reporting level
}

func hashTypeSeq(types []int32) string {
	h := sha1.New()
	b := make([]byte, len(types))
	for i, t := range types {
		b[i] = byte(t)
	}
	h.Write(b)
	return string(h.Sum(nil))
}

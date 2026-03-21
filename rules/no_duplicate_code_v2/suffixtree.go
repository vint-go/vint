package no_duplicate_code_v2

// Suffix tree implementation using Ukkonen's algorithm on []int32.
//
// All state and transition data is stored in flat []int32 arrays indexed
// by integer IDs. There are ZERO pointer-containing fields, so the GC's
// mark phase never needs to scan any of these arrays. This eliminates
// the massive GC overhead (54% of CPU) from the previous pointer-based
// implementation.

import "math"

const stInfinity = int32(math.MaxInt32)

// stree is a suffix tree operating on a flat []int32 data slice.
// States and transitions are referenced by int32 indices into parallel arrays.
type stree struct {
	data []int32

	// State storage (indexed by state ID).
	sLink       []int32 // suffix link state index (-1 = none)
	sTransHead  []int32 // first transition index (-1 = none, or start offset after compact)
	sTransCount []int32 // number of transitions (only valid after compact)
	numStates   int32

	// Transition storage (indexed by transition ID).
	tStart []int32 // edge start position in data
	tEnd   []int32 // edge end position in data
	tState []int32 // target state index
	tNext  []int32 // next sibling transition for same source state (-1 = last; nil after compact)
	numTrans int32

	root, aux int32

	// active point
	s          int32
	start, end int32
}

func newSTree(capacity int) *stree {
	// Ukkonen's creates at most 2n explicit states and 2n transitions.
	est := capacity*2 + 4
	t := &stree{
		data:       make([]int32, 0, capacity+1),
		sLink:      make([]int32, 0, est),
		sTransHead: make([]int32, 0, est),
		tStart:     make([]int32, 0, est),
		tEnd:       make([]int32, 0, est),
		tState:     make([]int32, 0, est),
		tNext:      make([]int32, 0, est),
	}
	t.root = t.allocState()
	t.aux = t.allocState()
	t.sLink[t.root] = t.aux
	t.s = t.root
	return t
}

func (t *stree) allocState() int32 {
	id := t.numStates
	t.numStates++
	t.sLink = append(t.sLink, -1)
	t.sTransHead = append(t.sTransHead, -1)
	return id
}

// addTran adds a transition from source state to target state.
// Returns the transition ID.
func (t *stree) addTran(state, start, end, target int32) int32 {
	id := t.numTrans
	t.numTrans++
	t.tStart = append(t.tStart, start)
	t.tEnd = append(t.tEnd, end)
	t.tState = append(t.tState, target)
	t.tNext = append(t.tNext, t.sTransHead[state])
	t.sTransHead[state] = id
	return id
}

// findTran finds the transition from state whose first data element matches c.
// Returns transition index or -1 if not found.
func (t *stree) findTran(state, c int32) int32 {
	for tr := t.sTransHead[state]; tr >= 0; tr = t.tNext[tr] {
		if t.data[t.tStart[tr]] == c {
			return tr
		}
	}
	return -1
}

func (t *stree) fork(state, i int32) int32 {
	r := t.allocState()
	t.addTran(state, i, stInfinity, r)
	return r
}

// hasTransitions returns true if the state has any outgoing transitions.
// Works in both linked-list mode (before compact) and contiguous mode (after compact).
func (t *stree) hasTransitions(state int32) bool {
	if t.sTransCount != nil {
		return t.sTransCount[state] > 0
	}
	return t.sTransHead[state] >= 0
}

// update adds tokens and extends the suffix tree incrementally.
func (t *stree) update(data []int32) {
	t.data = append(t.data, data...)
	for range data {
		t.extend()
		t.s, t.start = t.canonize(t.s, t.start, t.end)
		t.end++
	}
}

// extend transforms T(n) to T(n+1).
func (t *stree) extend() {
	oldr := t.root
	s := t.s
	start, end := t.start, t.end
	var r int32
	for {
		var endPoint bool
		r, endPoint = t.testAndSplit(s, start, end-1)
		if endPoint {
			break
		}
		t.fork(r, end)
		if oldr != t.root {
			t.sLink[oldr] = r
		}
		oldr = r
		s, start = t.canonize(t.sLink[s], start, end-1)
	}
	if oldr != t.root {
		t.sLink[oldr] = r
	}
	t.s = s
	t.start = start
}

func (t *stree) testAndSplit(s, start, end int32) (int32, bool) {
	c := t.data[t.end]
	if start <= end {
		tr := t.findTran(s, t.data[start])
		splitPoint := t.tStart[tr] + end - start + 1
		if t.data[splitPoint] == c {
			return s, true
		}
		ns := t.allocState()
		t.addTran(ns, splitPoint, t.tEnd[tr], t.tState[tr])
		t.tEnd[tr] = splitPoint - 1
		t.tState[tr] = ns
		return ns, false
	}
	if s == t.aux || t.findTran(s, c) >= 0 {
		return s, true
	}
	return s, false
}

func (t *stree) canonize(s, start, end int32) (int32, int32) {
	if s == t.aux {
		s, start = t.root, start+1
	}
	if start > end {
		return s, start
	}
	var tr int32
	for {
		if start <= end {
			tr = t.findTran(s, t.data[start])
		}
		if t.tEnd[tr]-t.tStart[tr] > end-start {
			break
		}
		start += t.tEnd[tr] - t.tStart[tr] + 1
		s = t.tState[tr]
	}
	return s, start
}

// compact reorganizes transitions into per-state contiguous blocks for
// cache-friendly access during tree walking. After this call, sTransHead[s]
// is the start offset and sTransCount[s] is the count. The tNext array is
// freed. Do not call update() after compact().
func (t *stree) compact() {
	n := t.numTrans
	newStart := make([]int32, n)
	newEnd := make([]int32, n)
	newState := make([]int32, n)

	// Count transitions per state.
	counts := make([]int32, t.numStates)
	for s := int32(0); s < t.numStates; s++ {
		for tr := t.sTransHead[s]; tr >= 0; tr = t.tNext[tr] {
			counts[s]++
		}
	}

	// Compute start offsets (prefix sum).
	offsets := make([]int32, t.numStates)
	off := int32(0)
	for s := int32(0); s < t.numStates; s++ {
		offsets[s] = off
		off += counts[s]
	}

	// Copy transitions into contiguous blocks.
	idx := make([]int32, t.numStates)
	copy(idx, offsets)
	for s := int32(0); s < t.numStates; s++ {
		for tr := t.sTransHead[s]; tr >= 0; tr = t.tNext[tr] {
			i := idx[s]
			newStart[i] = t.tStart[tr]
			newEnd[i] = t.tEnd[tr]
			newState[i] = t.tState[tr]
			idx[s]++
		}
	}

	t.tStart = newStart
	t.tEnd = newEnd
	t.tState = newState
	t.tNext = nil // no longer needed
	t.sTransHead = offsets
	t.sTransCount = counts
}

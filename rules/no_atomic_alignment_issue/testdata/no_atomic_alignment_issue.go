package fixtures

import "sync/atomic"

// --- Invalid examples ---

// Counter has a misaligned int64 field due to preceding bool (1 byte).
type Counter struct {
	flag  bool
	count int64 // misaligned on 32-bit
}

func (c *Counter) Increment() {
	atomic.AddInt64(&c.count, 1) // MATCH /address of non 64-bit aligned field .count passed to atomic function/
}

func (c *Counter) Load() int64 {
	return atomic.LoadInt64(&c.count) // MATCH /address of non 64-bit aligned field .count passed to atomic function/
}

// MisalignedUint64 has a misaligned uint64 field.
type MisalignedUint64 struct {
	x int32
	y uint64 // misaligned: preceded by 4-byte int32, offset = 4
}

func misalignedUint64Add(m *MisalignedUint64) {
	atomic.AddUint64(&m.y, 1) // MATCH /address of non 64-bit aligned field .y passed to atomic function/
}

// MisalignedAfterByte has a misaligned int64 after a byte field.
type MisalignedAfterByte struct {
	a byte
	b int64 // misaligned: preceded by 1-byte field
}

func misalignedStore(m *MisalignedAfterByte) {
	atomic.StoreInt64(&m.b, 42) // MATCH /address of non 64-bit aligned field .b passed to atomic function/
}

// --- Valid examples ---

// AlignedCounter has count as the first field, so it is properly aligned.
type AlignedCounter struct {
	count int64 // first field, offset = 0, always aligned
	flag  bool
}

func (c *AlignedCounter) Increment() {
	atomic.AddInt64(&c.count, 1) // safe: count is properly aligned
}

// AlignedAfterPadding has proper alignment because two int32s sum to 8 bytes.
type AlignedAfterPadding struct {
	a int32
	b int32
	c int64 // offset = 8, properly aligned
}

func alignedPaddingLoad(m *AlignedAfterPadding) int64 {
	return atomic.LoadInt64(&m.c) // safe: properly aligned
}

// AlignedAllInt64 has all int64 fields.
type AlignedAllInt64 struct {
	x int64
	y int64
}

func alignedAllInt64(m *AlignedAllInt64) {
	atomic.AddInt64(&m.y, 1) // safe: y is at offset 8
}

// Non-struct atomic usage (local variable) should not trigger.
func localVariable() {
	var x int64
	atomic.AddInt64(&x, 1)
}

// 32-bit operations should not trigger.
func atomic32Bit() {
	var x int32
	atomic.AddInt32(&x, 1)
}

package fixtures

// Invalid: unexported constant never referenced anywhere.

const maxRetries = 5 // MATCH /const maxRetries is unused/

func Process() error {
	return nil
}

// Invalid: both unexported constants in a block are unused (no iota, neither is referenced).

const (
	internalTimeout = 30  // MATCH /const internalTimeout is unused/
	internalBuffer  = 256 // MATCH /const internalBuffer is unused/
)

func Run() {}

// Valid: constant is referenced by a function.

const maxConnections = 10

func GetMaxConnections() int {
	return maxConnections
}

// Valid: exported constants are always considered used.

const Version = "2.0.0"

// Valid: if one constant in a block is used, all are considered used (iota group rule).

const (
	StatusPending  = iota
	StatusActive
	StatusInactive
	StatusArchived
)

func DefaultStatus() int {
	return StatusPending
}

// Valid: blank identifier is always considered used.

const _ = "blank identifier is always considered used"

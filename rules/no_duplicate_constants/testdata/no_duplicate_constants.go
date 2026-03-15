package fixtures

// Invalid: duplicate string constants
const (
	StatusActive = "active"
	UserActive   = "active" // MATCH /duplicate constant value "active": UserActive has the same value as StatusActive/
)

// Invalid: duplicate integer constants
const (
	DefaultTimeout = 30
	RequestTimeout = 30 // MATCH /duplicate constant value 30: RequestTimeout has the same value as DefaultTimeout/
	MaxRetryWait   = 30 // MATCH /duplicate constant value 30: MaxRetryWait has the same value as DefaultTimeout/
)

// Invalid: duplicate error string constants
const (
	ErrNotFound = "not found"
	ErrMissing  = "not found" // MATCH /duplicate constant value "not found": ErrMissing has the same value as ErrNotFound/
)

// Valid: all different values
const (
	StatusRunning  = "running"
	StatusInactive = "inactive"
	StatusPending  = "pending"
)

// Valid: referencing another constant (not a literal)
const (
	BaseTimeout    = 60
	RequestTimeout2 = BaseTimeout
)

// Valid: conceptually different constants with distinct values
const (
	MaxRetries      = 3
	DefaultTimeout3 = 45
	BufferSize      = 1024
)

// Valid: iota constants
const (
	A = iota
	B
	C
)

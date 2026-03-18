package fixtures

import "sync"

// Invalid: embedded sync.Mutex exposes Lock/Unlock as public methods.

type BadEmbeddedMutex struct {
	sync.Mutex // MATCH /embedded sync.Mutex exposes Lock and Unlock methods to the public API; use an unexported field instead/
	Data string
}

type BadEmbeddedRWMutex struct {
	sync.RWMutex // MATCH /embedded sync.RWMutex exposes Lock and Unlock methods to the public API; use an unexported field instead/
	Data string
}

// Invalid: exported named field of type sync.Mutex.

type BadExportedField struct {
	Mu   sync.Mutex // MATCH /exported field Mu of type sync.Mutex exposes Lock and Unlock methods to the public API; use an unexported field instead/
	Data string
}

type BadExportedRWField struct {
	Mu   sync.RWMutex // MATCH /exported field Mu of type sync.RWMutex exposes Lock and Unlock methods to the public API; use an unexported field instead/
	Data string
}

// Valid: unexported named field keeps Lock/Unlock private.

type GoodUnexportedMutex struct {
	mu   sync.Mutex
	Data string
}

type GoodUnexportedRWMutex struct {
	mu   sync.RWMutex
	Data string
}

// Valid: not a sync type.

type GoodNoSync struct {
	Name string
	Data int
}

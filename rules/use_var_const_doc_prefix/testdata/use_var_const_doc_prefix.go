package fixtures

// The default timeout for HTTP requests.
var DefaultTimeout = 30 // MATCH /comment on exported variable DefaultTimeout should be of the form "DefaultTimeout ..."/

// DefaultPort is the default port for the server.
var DefaultPort = 8080

// unexported variables should not be flagged
var unexportedVar = "hello"

// MaxRetries is the maximum number of retries.
var MaxRetries = 3

// A global flag for debug mode.
var DebugMode = false // MATCH /comment on exported variable DebugMode should be of the form "DebugMode ..."/

var NoDocVar = "no doc"

// The maximum number of connections allowed.
const MaxConnections = 100 // MATCH /comment on exported constant MaxConnections should be of the form "MaxConnections ..."/

// DefaultName is the default name.
const DefaultName = "default"

// unexported constants should not be flagged
const unexportedConst = 42

// Timeout is the default timeout duration.
const Timeout = 60

// A constant for pi approximation.
const Pi = 3.14 // MATCH /comment on exported constant Pi should be of the form "Pi ..."/

const NoDocConst = "no doc"

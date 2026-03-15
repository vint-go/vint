package fixtures

// Invalid: unexported package-level variable never read.

var unusedConfig = map[string]string{ // MATCH /var unusedConfig is unused/
	"key": "value",
}

var counter int // MATCH /var counter is unused/

func Increment() {
	// Writing to counter does not count as using it
	// (when post-statements-are-reads is false)
	counter++
}

var (
	Active   = true
	logLevel = "debug" // MATCH /var logLevel is unused/
)

func Status() bool {
	return Active
}

// Valid: variable is read by a function.

var config2 = map[string]string{
	"key": "value",
}

func GetConfig() map[string]string {
	return config2
}

// Valid: exported variables are always considered used.

var Version = "1.0.0"

// Valid: blank identifier is always used.

var _ = func() {}

// Valid: variable is read in init and then used.

var cache map[string]string

func init() {
	cache = make(map[string]string)
}

func Get(key string) string {
	return cache[key]
}

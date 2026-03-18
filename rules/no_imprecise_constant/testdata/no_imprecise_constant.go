package fixtures

import "math"

// Invalid: excessive precision in floating-point constants

const pi = 3.14159265358979323846264338327950288419716939937510 // MATCH /floating-point constant has excessive precision; use math constants or let Go handle the precision/

const e = 2.71828182845904523536028747135266249775724709369995 // MATCH /floating-point constant has excessive precision; use math constants or let Go handle the precision/

const sqrt2 = 1.41421356237309504880168872420969807856967187537694 // MATCH /floating-point constant has excessive precision; use math constants or let Go handle the precision/

const ln2 = 0.693147180559945309417232121458176568075500134360255 // MATCH /floating-point constant has excessive precision; use math constants or let Go handle the precision/

const (
	phi    = 1.61803398874989484820458683436563811772030917980576 // MATCH /floating-point constant has excessive precision; use math constants or let Go handle the precision/
	oneDiv3 = 0.333333333333333333333333333333333333 // MATCH /floating-point constant has excessive precision; use math constants or let Go handle the precision/
)

// Valid: reasonable precision (17 or fewer significant digits)
const piShort = 3.1415926535897932

const eShort = 2.7182818284590452

const small = 0.001

const one = 1.0

const half = 0.5

const largeRound = 1000000.0

// Valid: using math constants
var piVal = math.Pi

// Valid: integer constants are not affected
const maxInt = 1000000000000000000

// Valid: string constants are not affected
const greeting = "hello"

// Valid: exactly 17 significant digits
const exact17 = 3.1415926535897932

// Valid: short float
const shortFloat = 1.23e10

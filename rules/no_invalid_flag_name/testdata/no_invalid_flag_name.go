package fixtures

import "flag"

func invalidFlagNames() {
	flag.Bool("-verbose", false, "enable verbose mode")      // MATCH /flag name "-verbose" should not start with a hyphen/
	flag.String("output=file", "", "output file path")       // MATCH /flag name "output=file" should not contain '='/
	flag.Int("", 0, "some value")                            // MATCH /flag name is empty/
	flag.Duration("wait time", 0, "how long to wait")        // MATCH /flag name "wait time" should not contain whitespace/
	flag.Float64("-rate", 0.0, "the rate")                   // MATCH /flag name "-rate" should not start with a hyphen/
	flag.Int64("count=max", 0, "max count")                  // MATCH /flag name "count=max" should not contain '='/
	flag.Uint("", 0, "unsigned value")                       // MATCH /flag name is empty/
	flag.Uint64("-big", 0, "big number")                     // MATCH /flag name "-big" should not start with a hyphen/
	flag.Bool("", false, "empty bool flag")                  // MATCH /flag name is empty/
}

func invalidFlagNameVar() {
	var b bool
	flag.BoolVar(&b, "-debug", false, "debug mode")          // MATCH /flag name "-debug" should not start with a hyphen/
	var s string
	flag.StringVar(&s, "out=path", "", "output path")        // MATCH /flag name "out=path" should not contain '='/
	var i int
	flag.IntVar(&i, "", 0, "some int")                       // MATCH /flag name is empty/
}

func validFlagNames() {
	flag.Bool("verbose", false, "enable verbose mode")
	flag.String("output", "", "output file path")
	flag.Int("count", 0, "the count")
	flag.Duration("timeout", 0, "the timeout")
	flag.Float64("rate", 0.0, "the rate")
	flag.Int64("big", 0, "a big number")
	flag.Uint("unsigned", 0, "an unsigned value")
	flag.Uint64("huge", 0, "a huge number")
}

func validFlagNameVar() {
	var b bool
	flag.BoolVar(&b, "debug", false, "debug mode")
	var s string
	flag.StringVar(&s, "output", "", "output path")
	var i int
	flag.IntVar(&i, "count", 0, "some int")
}

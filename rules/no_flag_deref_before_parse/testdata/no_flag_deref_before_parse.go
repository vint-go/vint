package fixtures

import "flag"

func badUsage() {
	verbose := *flag.Bool("verbose", false, "enable verbose mode") // MATCH /immediate dereference of flag.Bool result before flag.Parse/
	_ = verbose
	flag.Parse()
}

func badUsageString() {
	name := *flag.String("name", "", "the name") // MATCH /immediate dereference of flag.String result before flag.Parse/
	_ = name
	flag.Parse()
}

func badUsageInt() {
	count := *flag.Int("count", 0, "the count") // MATCH /immediate dereference of flag.Int result before flag.Parse/
	_ = count
	flag.Parse()
}

func badUsageInt64() {
	big := *flag.Int64("big", 0, "a big number") // MATCH /immediate dereference of flag.Int64 result before flag.Parse/
	_ = big
	flag.Parse()
}

func badUsageUint() {
	u := *flag.Uint("u", 0, "unsigned") // MATCH /immediate dereference of flag.Uint result before flag.Parse/
	_ = u
	flag.Parse()
}

func badUsageUint64() {
	u64 := *flag.Uint64("u64", 0, "unsigned 64") // MATCH /immediate dereference of flag.Uint64 result before flag.Parse/
	_ = u64
	flag.Parse()
}

func badUsageFloat64() {
	f := *flag.Float64("f", 0.0, "a float") // MATCH /immediate dereference of flag.Float64 result before flag.Parse/
	_ = f
	flag.Parse()
}

func badUsageDuration() {
	d := *flag.Duration("d", 0, "a duration") // MATCH /immediate dereference of flag.Duration result before flag.Parse/
	_ = d
	flag.Parse()
}

func goodUsage() {
	verbose := flag.Bool("verbose", false, "enable verbose mode")
	flag.Parse()
	if *verbose {
		// use after Parse
	}
}

func goodUsageString() {
	name := flag.String("name", "", "the name")
	flag.Parse()
	_ = *name
}

func goodUsageNoDeref() {
	verbose := flag.Bool("verbose", false, "enable verbose mode")
	_ = verbose
}

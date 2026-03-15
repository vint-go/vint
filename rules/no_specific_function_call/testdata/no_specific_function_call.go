package fixtures

func bad() {
	println("hello") // MATCH /call of println(...)/
}

func alsobad() {
	println() // MATCH /call of println(...)/
}

func good() {
	example("hello") // not reported: different function name when configured with "println"
}

func goodQualified() {
	// Qualified calls like fmt.Println should not match
	_ = "fmt.Println"
}

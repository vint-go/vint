package fixtures

// Test that custom pattern "(?i)(my_custom_secret|my_custom_token)" is used

// Invalid: matches custom pattern
func customSecretAssignment() {
	my_custom_secret := "some_value_here" // MATCH /hardcoded credential: avoid embedding credentials in source code/
	_ = my_custom_secret
}

// Invalid: matches custom pattern
func customTokenAssignment() {
	my_custom_token := "another_value" // MATCH /hardcoded credential: avoid embedding credentials in source code/
	_ = my_custom_token
}

// Valid: default pattern "password" should NOT match with custom config
func defaultPatternNotMatching() {
	password := "supersecretpassword123"
	_ = password
}

// Valid: non-matching variable
func normalVariable() {
	name := "hello world"
	_ = name
}

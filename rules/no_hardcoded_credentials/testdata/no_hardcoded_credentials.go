package fixtures

import "os"

// Invalid: hardcoded password in variable assignment
func hardcodedPassword() {
	password := "supersecretpassword123" // MATCH /hardcoded credential: avoid embedding credentials in source code/
	_ = password
}

// Invalid: hardcoded API key in constant declaration
func hardcodedAPIKey() {
	const apiKey = "AKIAIOSFODNN7EXAMPLE" // MATCH /hardcoded credential: avoid embedding credentials in source code/
	_ = apiKey
}

// Invalid: hardcoded token in struct literal
type Config struct {
	Token    string
	Username string
}

func hardcodedTokenInStruct() {
	config := Config{
		Token: "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // MATCH /hardcoded credential: avoid embedding credentials in struct literals/
	}
	_ = config
}

// Invalid: password comparison with hardcoded value
func hardcodedComparison(userPassword string) {
	if userPassword == "admin123" { // MATCH /hardcoded credential: avoid comparing credentials against hardcoded values/
		grantAccess()
	}
}

// Valid: reading password from environment variable
func envPassword() {
	password := os.Getenv("DB_PASSWORD")
	_ = password
}

// Valid: using a function call to get secret
func functionCallSecret() {
	apiKey := getSecret("api-key")
	_ = apiKey
}

// Valid: non-credential variable with string value
func normalVariable() {
	name := "hello world"
	_ = name
}

// Valid: credential variable with empty string
func emptyCredential() {
	password := ""
	_ = password
}

// Valid: comparing non-credential variable
func normalComparison(status string) {
	if status == "active" {
		grantAccess()
	}
}

func grantAccess() {}

func getSecret(name string) string {
	return ""
}

package fixtures

const apiPrefix = "api/"

// Invalid: The literal "api/users" matches the evaluated constant expression
// apiPrefix + "users", so it appears 3 times total.
func getEndpoint1() string {
	return apiPrefix + "users" // MATCH /string literal "api/users" appears 3 times, consider extracting it into a named constant/
}

func getEndpoint2() string {
	return "api/users"
}

func getEndpoint3() string {
	return "api/users"
}

const baseURL = "https://example.com/"
const fullURL = baseURL + "api" // MATCH /string literal "https://example.com/api" appears 4 times, consider extracting it into a named constant/

// Invalid: The resolved constant fullURL equals "https://example.com/api",
// which also appears as a literal and as baseURL + "api".
func getURL1() string {
	return fullURL
}

func getURL2() string {
	return "https://example.com/api"
}

func getURL3() string {
	return baseURL + "api"
}

func getURL4() string {
	return "https://example.com/api"
}

// Valid: The concatenation result "api/health" only appears twice (below threshold of 3).
func healthCheck1() string {
	return apiPrefix + "health"
}

func healthCheck2() string {
	return "api/health"
}

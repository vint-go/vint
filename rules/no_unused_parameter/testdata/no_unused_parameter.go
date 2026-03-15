package fixtures

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Invalid: parameter "name" is never used in the function body.
func greet(name string) string { // MATCH /parameter 'name' seems to be unused, consider removing or renaming it as _/
	return "hello"
}

// Invalid: parameter "count" is never referenced.
func processItems(items []string, count int) error { // MATCH /parameter 'count' seems to be unused, consider removing or renaming it as _/
	for _, item := range items {
		fmt.Println(item)
	}
	return nil
}

// Invalid: parameters "ctx" and "logger" are unused.
func fetchData(
	ctx context.Context, // MATCH /parameter 'ctx' seems to be unused, consider removing or renaming it as _/
	url string,
	logger *log.Logger, // MATCH /parameter 'logger' seems to be unused, consider removing or renaming it as _/
) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// Valid: all parameters are used in the function body.
func greetUsed(name string) string {
	return "hello, " + name
}

// Valid: all parameters are referenced.
func processItemsUsed(items []string, count int) error {
	if len(items) != count {
		return fmt.Errorf("expected %d items, got %d", count, len(items))
	}
	for _, item := range items {
		fmt.Println(item)
	}
	return nil
}

// Valid: using the blank identifier to explicitly mark an unused parameter.
func handler(_ http.ResponseWriter, r *http.Request) {
	log.Println("received request:", r.URL.Path)
}

// Valid: all parameters are used, including the logger.
func fetchDataUsed(ctx context.Context, url string, logger *log.Logger) ([]byte, error) {
	logger.Printf("fetching %s", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// Valid: exported function is not checked by default.
func ExportedFunc(unused string) string {
	return "hello"
}

// Valid: function with no parameters.
func noParams() {}

// Valid: function with only blank parameter.
func blankOnly(_ int) {}

// Valid: function with no body (prototype).
func protoFunc(x int)

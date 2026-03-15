package fixtures

import (
	"context"
	"fmt"
	"net/http"
)

// Invalid: calling (*http.Client).Get without a context
func badClientGet() {
	client := &http.Client{}
	resp, err := client.Get("https://example.com") // MATCH /(*http.Client).Get does not accept a context; use (*http.Client).Do with http.NewRequestWithContext instead/
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}

// Valid: using (*http.Client).Do with NewRequestWithContext
func goodClientDoWithContext() {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}

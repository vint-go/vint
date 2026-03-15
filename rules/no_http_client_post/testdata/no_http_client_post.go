package fixtures

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Invalid: calling (*http.Client).Post without a context
func badClientPost() {
	client := &http.Client{}
	body := strings.NewReader(`{"key": "value"}`)
	resp, err := client.Post("https://example.com/api", "application/json", body) // MATCH /(*http.Client).Post does not accept a context; use (*http.Client).Do with http.NewRequestWithContext instead/
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}

// Valid: using (*http.Client).Do with NewRequestWithContext
func goodClientDoWithContext() {
	ctx := context.Background()
	body := strings.NewReader(`{"key": "value"}`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://example.com/api", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}

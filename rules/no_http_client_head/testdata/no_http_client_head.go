package fixtures

import (
	"context"
	"fmt"
	"net/http"
)

// Invalid: calling (*http.Client).Head without a context
func badClientHead() {
	client := &http.Client{}
	resp, err := client.Head("https://example.com") // MATCH /(*http.Client).Head does not accept a context; use (*http.Client).Do with http.NewRequestWithContext instead/
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}

// Valid: using (*http.Client).Do with NewRequestWithContext
func goodClientDoWithContext() {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, "https://example.com", nil)
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

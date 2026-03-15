package fixtures

import (
	"context"
	"fmt"
	"net/http"
)

func badHttpGet() {
	resp, err := http.Get("https://example.com") // MATCH /http.Get does not accept a context; use http.NewRequestWithContext and (*http.Client).Do instead/
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}

func goodHttpGetWithContext() {
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

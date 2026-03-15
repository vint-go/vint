package fixtures

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func badHttpPost() {
	body := strings.NewReader(`{"key": "value"}`)
	resp, err := http.Post("https://example.com/api", "application/json", body) // MATCH /http.Post does not accept a context; use http.NewRequestWithContext and (*http.Client).Do instead/
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)
}

func goodHttpPostWithContext() {
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

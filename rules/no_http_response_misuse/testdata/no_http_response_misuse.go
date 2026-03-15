package fixtures

import (
	"io"
	"log"
	"net/http"
)

// Invalid: deferring Body.Close before checking error
func httpResponseMisuseGet() {
	resp, err := http.Get("https://example.com")
	defer resp.Body.Close() // MATCH /deferring Body.Close before checking error from HTTP call may cause nil pointer dereference/
	if err != nil {
		log.Fatal(err)
	}
	_, _ = io.ReadAll(resp.Body)
}

// Invalid: deferring Body.Close before checking error with Client.Do
func httpResponseMisuseClientDo(client *http.Client) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	resp, err := client.Do(req)
	defer resp.Body.Close() // MATCH /deferring Body.Close before checking error from HTTP call may cause nil pointer dereference/
	if err != nil {
		log.Fatal(err)
	}
	_ = resp.StatusCode
}

// Valid: check error before deferring Body.Close
func httpResponseCorrectGet() {
	resp, err := http.Get("https://example.com")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
}

// Valid: check error before deferring Body.Close with Client.Do
func httpResponseCorrectClientDo(client *http.Client) {
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	_ = resp.StatusCode
}

// Valid: blank identifier for response
func httpResponseBlankIdentifier() {
	_, err := http.Get("https://example.com")
	if err != nil {
		log.Fatal(err)
	}
}

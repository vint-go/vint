package fixtures

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

func noUncheckedErrorOpenIgnored() {
	os.Open("file.txt") // MATCH /unchecked error in call to os.Open/
}

func noUncheckedErrorBlankIdentifier() {
	f, _ := os.Open("file.txt")
	defer f.Close() // MATCH /unchecked error in call to os.File.Close/
}

func noUncheckedErrorUnmarshalIgnored() {
	var data map[string]interface{}
	json.Unmarshal([]byte(`{"key":"value"}`), &data) // MATCH /unchecked error in call to encoding/json.Unmarshal/
}

func noUncheckedErrorListenAndServeIgnored() {
	http.ListenAndServe(":8080", nil) // MATCH /unchecked error in call to net/http.ListenAndServe/
}

func noUncheckedErrorOpenChecked() {
	f, err := os.Open("file.txt")
	if err != nil {
		return
	}
	_ = f
}

func noUncheckedErrorUnmarshalChecked() {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(`{"key":"value"}`), &data)
	if err != nil {
		return
	}
}

func noUncheckedErrorListenAndServeChecked() {
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return
	}
}

func noUncheckedErrorFmtPrintlnExcluded() {
	fmt.Println("Hello, world!")
}

func noUncheckedErrorCloseIgnored() {
	f, err := os.Open("file.txt")
	if err != nil {
		return
	}
	defer f.Close() // MATCH /unchecked error in call to os.File.Close/
}

func noUncheckedErrorInsideGoroutine() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.Get("http://example.com") // MATCH /unchecked error in call to net/http.Get/
	}()
	wg.Wait()
}

func noUncheckedErrorInsideDeferFuncLit() {
	defer func() {
		os.Remove("/tmp/test") // MATCH /unchecked error in call to os.Remove/
	}()
}

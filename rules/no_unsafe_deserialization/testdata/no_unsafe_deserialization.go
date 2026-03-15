package fixtures

import (
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

// Invalid: gob.NewDecoder from HTTP request body (untrusted source)
func handlerGob(w http.ResponseWriter, r *http.Request) {
	decoder := gob.NewDecoder(r.Body) // MATCH /unsafe deserialization: untrusted data passed to gob.NewDecoder/
	var data interface{}
	decoder.Decode(&data)
}

// Invalid: xml.Unmarshal of untrusted data from HTTP request body
func handlerXML(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var result interface{}
	xml.Unmarshal(body, &result) // MATCH /unsafe deserialization: untrusted data passed to xml.Unmarshal/
}

// Valid: json.NewDecoder into strongly-typed struct with validation
type UserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func handlerJSONSafe(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var input UserInput
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if input.Name == "" || input.Email == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}
}

// Valid: json.NewDecoder with concrete type
type AppConfig struct {
	Host string
	Port int
}

func handlerJSONConfig(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var config AppConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
}

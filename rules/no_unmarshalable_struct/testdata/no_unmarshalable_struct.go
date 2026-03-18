package fixtures

import (
	"encoding/json"
	"encoding/xml"
)

// Invalid: struct with no exported fields being marshaled
type config struct {
	name  string `json:"name"`
	value int    `json:"value"`
}

func badMarshalNoExported() {
	c := config{name: "test", value: 42}
	json.Marshal(c) // MATCH /struct config has no exported fields and will marshal as an empty object/
}

// Invalid: struct with no exported fields using xml.Marshal
type xmlConfig struct {
	host string `xml:"host"`
	port int    `xml:"port"`
}

func badXmlMarshalNoExported() {
	c := xmlConfig{host: "localhost", port: 8080}
	xml.Marshal(c) // MATCH /struct xmlConfig has no exported fields and will marshal as an empty object/
}

// Invalid: struct with no exported fields using json.MarshalIndent
func badMarshalIndentNoExported() {
	c := config{name: "test", value: 42}
	json.MarshalIndent(c, "", "  ") // MATCH /struct config has no exported fields and will marshal as an empty object/
}

// Invalid: struct with no exported fields using xml.MarshalIndent
func badXmlMarshalIndentNoExported() {
	c := xmlConfig{host: "localhost", port: 8080}
	xml.MarshalIndent(c, "", "  ") // MATCH /struct xmlConfig has no exported fields and will marshal as an empty object/
}

// Invalid: struct with no exported fields using json.Encoder.Encode
func badJsonEncodeNoExported() {
	c := config{name: "test", value: 42}
	enc := json.NewEncoder(nil)
	enc.Encode(c) // MATCH /struct config has no exported fields and will marshal as an empty object/
}

// Invalid: struct with no exported fields using xml.Encoder.Encode
func badXmlEncodeNoExported() {
	c := xmlConfig{host: "localhost", port: 8080}
	enc := xml.NewEncoder(nil)
	enc.Encode(c) // MATCH /struct xmlConfig has no exported fields and will marshal as an empty object/
}

// Invalid: pointer to struct with no exported fields
func badMarshalPointerNoExported() {
	c := &config{name: "test", value: 42}
	json.Marshal(c) // MATCH /struct config has no exported fields and will marshal as an empty object/
}

// Valid: struct with exported fields
type Config struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func goodMarshalExported() {
	c := Config{Name: "test", Value: 42}
	json.Marshal(c)
}

// Valid: struct with some exported fields
type MixedConfig struct {
	Name     string `json:"name"`
	internal string
}

func goodMarshalMixed() {
	c := MixedConfig{Name: "test"}
	_ = c.internal
	json.Marshal(c)
}

// Valid: marshaling a non-struct type
func goodMarshalString() {
	s := "hello"
	json.Marshal(s)
}

// Valid: marshaling a map
func goodMarshalMap() {
	m := map[string]int{"a": 1}
	json.Marshal(m)
}

// Valid: marshaling a slice
func goodMarshalSlice() {
	s := []int{1, 2, 3}
	json.Marshal(s)
}

// Valid: struct with custom MarshalJSON method
type customMarshaler struct {
	data string
}

func (c customMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(`"custom"`), nil
}

func goodMarshalCustom() {
	c := customMarshaler{data: "test"}
	json.Marshal(c)
}

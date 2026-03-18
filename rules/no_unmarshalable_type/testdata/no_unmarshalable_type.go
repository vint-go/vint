package fixtures

import (
	"encoding/json"
	"encoding/xml"
)

// Invalid: marshaling a channel directly
func badMarshalChannel() {
	ch := make(chan int)
	json.Marshal(ch) // MATCH /cannot marshal channels or functions/
}

// Invalid: marshaling a function directly
func badMarshalFunc() {
	fn := func() {}
	json.Marshal(fn) // MATCH /cannot marshal channels or functions/
}

// Invalid: marshaling a channel with xml.Marshal
func badXmlMarshalChannel() {
	ch := make(chan string)
	xml.Marshal(ch) // MATCH /cannot marshal channels or functions/
}

// Invalid: marshaling a channel with json.MarshalIndent
func badMarshalIndentChannel() {
	ch := make(chan int)
	json.MarshalIndent(ch, "", "  ") // MATCH /cannot marshal channels or functions/
}

// Invalid: marshaling a channel with xml.MarshalIndent
func badXmlMarshalIndentChannel() {
	ch := make(chan int)
	xml.MarshalIndent(ch, "", "  ") // MATCH /cannot marshal channels or functions/
}

// Invalid: marshaling a channel with json.Encoder.Encode
func badJsonEncodeChannel() {
	ch := make(chan int)
	enc := json.NewEncoder(nil)
	enc.Encode(ch) // MATCH /cannot marshal channels or functions/
}

// Invalid: marshaling a function with xml.Encoder.Encode
func badXmlEncodeFunc() {
	fn := func() {}
	enc := xml.NewEncoder(nil)
	enc.Encode(fn) // MATCH /cannot marshal channels or functions/
}

// Valid: marshaling a struct with normal fields
type Data struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func goodMarshalStruct() {
	d := Data{Name: "test", Value: 42}
	json.Marshal(d)
}

// Valid: marshaling a string
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

// Valid: marshaling an interface value
func goodMarshalInterface() {
	var v interface{} = "hello"
	json.Marshal(v)
}

// Valid: marshaling a pointer to struct
func goodMarshalPointer() {
	d := &Data{Name: "test", Value: 42}
	json.Marshal(d)
}

package fixtures

import (
	"encoding/json"
	"encoding/xml"
)

type User struct {
	Name string `json:"name"`
}

// Invalid: passing non-pointer to json.Unmarshal
func badJsonUnmarshal() {
	var u User
	data := []byte(`{"name":"Alice"}`)
	json.Unmarshal(data, u) // MATCH /call of unmarshal-like function with non-pointer argument/
}

// Invalid: passing non-pointer to xml.Unmarshal
func badXmlUnmarshal() {
	var u User
	data := []byte(`<User><Name>Alice</Name></User>`)
	xml.Unmarshal(data, u) // MATCH /call of unmarshal-like function with non-pointer argument/
}

// Invalid: passing non-pointer to json.Decoder.Decode
func badJsonDecode() {
	var u User
	dec := json.NewDecoder(nil)
	dec.Decode(u) // MATCH /call of unmarshal-like function with non-pointer argument/
}

// Invalid: passing non-pointer to xml.Decoder.Decode
func badXmlDecode() {
	var u User
	dec := xml.NewDecoder(nil)
	dec.Decode(u) // MATCH /call of unmarshal-like function with non-pointer argument/
}

// Valid: passing pointer to json.Unmarshal
func goodJsonUnmarshal() {
	var u User
	data := []byte(`{"name":"Alice"}`)
	json.Unmarshal(data, &u)
}

// Valid: passing pointer to xml.Unmarshal
func goodXmlUnmarshal() {
	var u User
	data := []byte(`<User><Name>Alice</Name></User>`)
	xml.Unmarshal(data, &u)
}

// Valid: passing pointer to json.Decoder.Decode
func goodJsonDecode() {
	var u User
	dec := json.NewDecoder(nil)
	dec.Decode(&u)
}

// Valid: passing interface value to json.Unmarshal
func goodJsonUnmarshalInterface() {
	data := []byte(`{"name":"Alice"}`)
	var v interface{}
	json.Unmarshal(data, &v)
}

// Valid: passing a map pointer (map is reference type but still passing pointer)
func goodJsonUnmarshalMapPtr() {
	data := []byte(`{"name":"Alice"}`)
	var m map[string]interface{}
	json.Unmarshal(data, &m)
}

package fixtures

// Invalid examples

type BadDuplicate struct {
	Name string `json:"name" json:"username"` // MATCH /malformed struct tag: duplicate tag key "json"/
}

type BadMalformed struct {
	Name string `json:"name` // MATCH /malformed struct tag: badly quoted value in tag/
}

type BadSpace struct {
	Name string `json: "name"` // MATCH /malformed struct tag: key "json" has space before value in tag/
}

// Valid examples

type GoodUser struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	Age   int    `json:"age" xml:"age"`
}

type GoodExcluded struct {
	Name string `json:"name" validate:"required"`
	ID   int    `json:"-"` // excluded from JSON
}

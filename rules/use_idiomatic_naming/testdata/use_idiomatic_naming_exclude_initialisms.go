package fixtures

// With excludeInitialisms: ["HTTP", "URL"]
// HTTP and URL are excluded from initialism checks, so Http/Url forms are accepted.

// Valid: Excluded initialisms are not enforced
type HttpClient struct{}

func GetUrlPath() string { return "" }

// Invalid: Non-excluded initialisms are still enforced
type XmlParser struct{} // MATCH /type XmlParser should be XMLParser/

var myId int // MATCH /var myId should be myID/

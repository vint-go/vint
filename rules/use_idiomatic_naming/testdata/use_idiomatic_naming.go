package fixtures

// Invalid: snake_case names
var my_variable int // MATCH /don't use underscores in Go names; var my_variable should be myVariable/

// Invalid: Incorrect casing for acronyms
type HttpClient struct{} // MATCH /type HttpClient should be HTTPClient/
type XmlParser struct{}  // MATCH /type XmlParser should be XMLParser/

func GetUrlPath() string { return "" } // MATCH /func GetUrlPath should be GetURLPath/

// Invalid: More acronym casing issues
func GetHtmlContent() string { return "" } // MATCH /func GetHtmlContent should be GetHTMLContent/

var myId int // MATCH /var myId should be myID/

type JsonResponse struct{} // MATCH /type JsonResponse should be JSONResponse/

// Valid: camelCase names
var myVariable int

// Valid: Correct casing for acronyms
type HTTPClient struct{}

type XMLParser struct{}

func GetURLPath() string { return "" }

func GetHTMLContent() string { return "" }

var myID int

type JSONResponse struct{}

// Valid: simple lowercase names
var simple int

// Valid: blank identifier
var _ int

// Valid: single letter names
var x int

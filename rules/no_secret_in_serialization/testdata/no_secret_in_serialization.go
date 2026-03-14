package fixtures

// Invalid: Password exposed via JSON serialization
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"` // MATCH /field "Password" with secret-like name is exposed via json serialization, use `json:"-"` to exclude it/
}

// Invalid: APIKey exposed via YAML serialization
type Config struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	APIKey string `yaml:"api_key"` // MATCH /field "APIKey" with secret-like name is exposed via yaml serialization, use `yaml:"-"` to exclude it/
}

// Invalid: Secret exposed via XML serialization
type XMLConfig struct {
	Name   string `xml:"name"`
	Secret string `xml:"secret"` // MATCH /field "Secret" with secret-like name is exposed via xml serialization, use `xml:"-"` to exclude it/
}

// Invalid: Token exposed via TOML serialization
type TOMLConfig struct {
	Name  string `toml:"name"`
	Token string `toml:"token"` // MATCH /field "Token" with secret-like name is exposed via toml serialization, use `toml:"-"` to exclude it/
}

// Invalid: AuthToken exposed via JSON with omitempty
type AuthConfig struct {
	AuthToken string `json:"auth_token,omitempty"` // MATCH /field "AuthToken" with secret-like name is exposed via json serialization, use `json:"-"` to exclude it/
}

// Invalid: AccessKey exposed via JSON
type AWSConfig struct {
	Region    string `json:"region"`
	AccessKey string `json:"access_key"` // MATCH /field "AccessKey" with secret-like name is exposed via json serialization, use `json:"-"` to exclude it/
}

// Invalid: PrivateKey exposed via JSON
type CryptoConfig struct {
	PrivateKey string `json:"private_key"` // MATCH /field "PrivateKey" with secret-like name is exposed via json serialization, use `json:"-"` to exclude it/
}

// Invalid: Credential exposed via JSON
type ServiceConfig struct {
	Credential string `json:"credential"` // MATCH /field "Credential" with secret-like name is exposed via json serialization, use `json:"-"` to exclude it/
}

// Valid: Password excluded from JSON output
type SafeUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

// Valid: Using separate response type without sensitive fields
type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Valid: Secret excluded from YAML output
type SafeConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	APIKey string `yaml:"-"`
}

// Valid: No serialization tags on secret field
type InternalUser struct {
	Name     string `json:"name"`
	Password string
}

// Valid: Non-secret field with serialization tags
type Product struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// Valid: Secret excluded from XML output
type SafeXMLConfig struct {
	Name   string `xml:"name"`
	Secret string `xml:"-"`
}

// Valid: Token excluded from TOML output
type SafeTOMLConfig struct {
	Name  string `toml:"name"`
	Token string `toml:"-"`
}

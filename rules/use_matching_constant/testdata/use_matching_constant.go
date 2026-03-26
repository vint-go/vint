package fixtures

const StatusActive = "active"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// Invalid: string literal "active" appears 3 times (meets min-occurrences threshold of 3)
// and matches existing constant StatusActive.

func IsActive(status string) bool {
	return status == "active" // MATCH /string literal "active" matches constant StatusActive, use the constant instead/
}

func IsActive2(status string) bool {
	return status == "active" // MATCH /string literal "active" matches constant StatusActive, use the constant instead/
}

func IsActive3(status string) bool {
	return status == "active" // MATCH /string literal "active" matches constant StatusActive, use the constant instead/
}

// Valid: "admin" appears only once as a non-map-key literal (below threshold of 3).

func CheckRole(role string) bool {
	if role == "admin" {
		return true
	}
	return false
}

// Valid: correctly referencing the named constants.

func IsActiveGood(status string) bool {
	return status == StatusActive
}

func CheckRoleGood(role string) bool {
	if role == RoleAdmin {
		return true
	}
	if role == RoleUser {
		return true
	}
	return false
}

// Valid: no constant exists for "pending", so a raw literal is fine.

func IsPending(status string) bool {
	return status == "pending"
}

// Valid: "user" used as map key should not be counted.

func MapKeyExample() map[string]string {
	return map[string]string{
		"user": "some value",
	}
}

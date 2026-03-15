package fixtures

const StatusActive = "active"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// Invalid: string literals matching existing constants.

func IsActive(status string) bool {
	return status == "active" // MATCH /string literal "active" matches constant StatusActive, use the constant instead/
}

func CheckRole(role string) bool {
	if role == "admin" { // MATCH /string literal "admin" matches constant RoleAdmin, use the constant instead/
		return true
	}
	if role == "user" { // MATCH /string literal "user" matches constant RoleUser, use the constant instead/
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

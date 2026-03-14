package fixtures

import "net/http"

// Invalid: Missing all security attributes (three failures)
func missingAll(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{ /* MATCH /http.Cookie should set Secure to true to prevent transmission over unencrypted connections/
		MATCH /http.Cookie should set HttpOnly to true to prevent JavaScript access/
		MATCH /http.Cookie should set SameSite to prevent CSRF attacks/ */
		Name:  "session",
		Value: "abc123",
	})
}

// Invalid: Secure and HttpOnly explicitly set to false (two failures)
func explicitFalse(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{ /* MATCH /http.Cookie should set Secure to true to prevent transmission over unencrypted connections/
		MATCH /http.Cookie should set HttpOnly to true to prevent JavaScript access/
		MATCH /http.Cookie should set SameSite to prevent CSRF attacks/ */
		Name:     "session",
		Value:    "abc123",
		Secure:   false,
		HttpOnly: false,
	})
}

// Invalid: Only Secure is set, missing HttpOnly and SameSite
func partialSecureOnly(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{ /* MATCH /http.Cookie should set HttpOnly to true to prevent JavaScript access/
		MATCH /http.Cookie should set SameSite to prevent CSRF attacks/ */
		Name:   "session",
		Value:  "abc123",
		Secure: true,
	})
}

// Invalid: Only HttpOnly is set, missing Secure and SameSite
func partialHttpOnlyOnly(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{ /* MATCH /http.Cookie should set Secure to true to prevent transmission over unencrypted connections/
		MATCH /http.Cookie should set SameSite to prevent CSRF attacks/ */
		Name:     "session",
		Value:    "abc123",
		HttpOnly: true,
	})
}

// Invalid: Has Secure and HttpOnly but missing SameSite
func missingSameSite(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{ // MATCH /http.Cookie should set SameSite to prevent CSRF attacks/
		Name:     "session",
		Value:    "abc123",
		Secure:   true,
		HttpOnly: true,
	})
}

// Invalid: Has all attributes but Secure is false
func secureFalse(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{ // MATCH /http.Cookie should set Secure to true to prevent transmission over unencrypted connections/
		Name:     "session",
		Value:    "abc123",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

// Invalid: Has all attributes but HttpOnly is false
func httpOnlyFalse(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{ // MATCH /http.Cookie should set HttpOnly to true to prevent JavaScript access/
		Name:     "session",
		Value:    "abc123",
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
	})
}

// Valid: All security attributes properly set
func allSecure(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "abc123",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

// Valid: Using SameSiteLaxMode is also acceptable
func allSecureLax(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "abc123",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// Valid: Cookie used directly (not via pointer)
func directCookie(w http.ResponseWriter) {
	c := http.Cookie{
		Name:     "session",
		Value:    "abc123",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	_ = c
}

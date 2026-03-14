package fixtures

import (
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
)

// Invalid: Direct taint from HTTP request to smtp.SendMail
func handlerSendMail(w http.ResponseWriter, r *http.Request) {
	to := r.FormValue("to")
	subject := r.FormValue("subject")
	body := r.FormValue("body")

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body)

	smtp.SendMail("smtp.example.com:587", nil, "from@example.com", []string{to}, msg) // MATCH /potential SMTP injection: tainted data from user input flows to smtp.SendMail/
}

// Invalid: Taint flows through fmt.Sprintf to smtp.SendMail
func handlerSprintf(w http.ResponseWriter, r *http.Request) {
	to := r.FormValue("to")
	subject := r.FormValue("subject")

	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\nHello", to, subject)
	smtp.SendMail("smtp.example.com:587", nil, "from@example.com", []string{to}, []byte(msg)) // MATCH /potential SMTP injection: tainted data from user input flows to smtp.SendMail/
}

// Invalid: Interprocedural - taint flows through helper function
func sendUnsafe(to string, msgBody []byte) {
	smtp.SendMail("smtp.example.com:587", nil, "from@example.com", []string{to}, msgBody)
}

func handlerInterprocedural(w http.ResponseWriter, r *http.Request) {
	to := r.FormValue("to")
	body := r.FormValue("body")
	msg := []byte("To: " + to + "\r\n\r\n" + body)
	sendUnsafe(to, msg) // MATCH /potential SMTP injection: tainted data from user input flows to smtp.SendMail/
}

// Invalid: Taint flows through string concatenation
func handlerConcat(w http.ResponseWriter, r *http.Request) {
	to := r.URL.Query().Get("to")
	msg := "To: " + to + "\r\n\r\nHello"
	smtp.SendMail("smtp.example.com:587", nil, "from@example.com", []string{to}, []byte(msg)) // MATCH /potential SMTP injection: tainted data from user input flows to smtp.SendMail/
}

// Valid: Sanitized input - CRLF removed before use
func handlerSanitized(w http.ResponseWriter, r *http.Request) {
	to := r.FormValue("to")
	subject := r.FormValue("subject")
	body := r.FormValue("body")

	cleanTo := strings.NewReplacer("\r", "", "\n", "").Replace(to)
	cleanSubject := strings.NewReplacer("\r", "", "\n", "").Replace(subject)

	msg := []byte("To: " + cleanTo + "\r\n" +
		"Subject: " + cleanSubject + "\r\n" +
		"\r\n" + body)

	smtp.SendMail("smtp.example.com:587", nil, "from@example.com", []string{cleanTo}, msg)
}

// Valid: No user input in SMTP call
func safeSendMail() {
	msg := []byte("To: admin@example.com\r\nSubject: Alert\r\n\r\nSystem alert")
	smtp.SendMail("smtp.example.com:587", nil, "from@example.com", []string{"admin@example.com"}, msg)
}

// Valid: Hardcoded constant values
func sendWelcomeEmail() {
	to := "user@example.com"
	msg := []byte("To: " + to + "\r\nSubject: Welcome\r\n\r\nWelcome to our service!")
	smtp.SendMail("smtp.example.com:587", nil, "from@example.com", []string{to}, msg)
}

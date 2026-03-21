package example

func getStatus() string {
	return "active" //nolint:goconst
}

func checkStatus(s string) bool {
	return s == "active"
}

func defaultStatus() string {
	return "active"
}

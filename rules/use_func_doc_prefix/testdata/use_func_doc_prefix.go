package fixtures

// This function processes data.
func ProcessData() {} // MATCH /comment on exported function ProcessData should be of the form "ProcessData ..."/

// ProcessDataCorrect processes the given data and returns results.
func ProcessDataCorrect() {}

// unexported functions should not be flagged
func unexportedFunc() {}

// Helper is a helper function that does something useful.
func Helper() {}

// Does something important.
func DoSomething() {} // MATCH /comment on exported function DoSomething should be of the form "DoSomething ..."/

func NoDocFunc() {}

// RunTask starts the given task.
func RunTask() {}

// Starts the server.
func StartServer() {} // MATCH /comment on exported function StartServer should be of the form "StartServer ..."/

// ValidateInput validates the input data.
func ValidateInput() {}

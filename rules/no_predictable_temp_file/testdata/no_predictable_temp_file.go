package fixtures

import "os"

func writeTemp() error {
	// Predictable temporary file path
	f, err := os.Create("/tmp/myapp.log") // MATCH /predictable temporary file path, use os.CreateTemp instead/
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("log data")
	return err
}

func writeTempData(data []byte) error {
	// Predictable path in temp directory
	return os.WriteFile("/tmp/myapp-cache.dat", data, 0600) // MATCH /predictable temporary file path, use os.CreateTemp instead/
}

func openTempFile() (*os.File, error) {
	// Predictable path via os.OpenFile
	return os.OpenFile("/tmp/myapp.lock", os.O_CREATE|os.O_WRONLY, 0600) // MATCH /predictable temporary file path, use os.CreateTemp instead/
}

func mkdirTemp() error {
	// Predictable path via os.Mkdir
	return os.Mkdir("/tmp/myapp-data", 0700) // MATCH /predictable temporary file path, use os.CreateTemp instead/
}

func mkdirAllTemp() error {
	// Predictable path via os.MkdirAll
	return os.MkdirAll("/tmp/myapp/subdir", 0700) // MATCH /predictable temporary file path, use os.CreateTemp instead/
}

func varTmpPath() error {
	// Predictable path in /var/tmp/
	return os.WriteFile("/var/tmp/myapp.dat", []byte("data"), 0600) // MATCH /predictable temporary file path, use os.CreateTemp instead/
}

// Valid: Using os.CreateTemp for unique file name
func writeTempSafe() error {
	f, err := os.CreateTemp("", "myapp-*.log")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("log data")
	return err
}

// Valid: Using os.CreateTemp with a specific directory
func writeTempDataSafe(data []byte) error {
	f, err := os.CreateTemp(os.TempDir(), "myapp-cache-*.dat")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

// Valid: Writing to non-temp directory
func writeNonTemp(data []byte) error {
	return os.WriteFile("/var/log/myapp.log", data, 0600)
}

// Valid: Creating file in current directory
func writeCurrentDir() error {
	return os.WriteFile("output.txt", []byte("data"), 0600)
}

// Valid: Non-string-literal argument
func writeDynamic(path string, data []byte) error {
	return os.WriteFile(path, data, 0600)
}

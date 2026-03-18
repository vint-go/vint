package fixtures

import "os"

func badWriteFile() {
	os.WriteFile("file.txt", []byte("data"), 644) // MATCH /file mode 644 is not in octal; did you mean 0644?/
}

func badMkdir() {
	os.Mkdir("dir", 755) // MATCH /file mode 755 is not in octal; did you mean 0755?/
}

func badMkdirAll() {
	os.MkdirAll("dir/sub", 700) // MATCH /file mode 700 is not in octal; did you mean 0700?/
}

func badOpenFile() {
	os.OpenFile("file.txt", os.O_CREATE, 666) // MATCH /file mode 666 is not in octal; did you mean 0666?/
}

func badChmod() {
	os.Chmod("file.txt", 444) // MATCH /file mode 444 is not in octal; did you mean 0444?/
}

// Valid examples below: these should not trigger failures

func goodWriteFileOctal() {
	os.WriteFile("file.txt", []byte("data"), 0644)
}

func goodMkdirOctal() {
	os.Mkdir("dir", 0755)
}

func goodMkdirAllOctal() {
	os.MkdirAll("dir/sub", 0700)
}

func goodOpenFileOctal() {
	os.OpenFile("file.txt", os.O_CREATE, 0666)
}

func goodChmodOctal() {
	os.Chmod("file.txt", 0444)
}

func goodWriteFileModernOctal() {
	os.WriteFile("file.txt", []byte("data"), 0o644)
}

func goodWriteFileZero() {
	os.WriteFile("file.txt", []byte("data"), 0)
}

func goodDecimalWithEightOrNine() {
	// 890 contains digits > 7, so it clearly wasn't intended as octal
	os.WriteFile("file.txt", []byte("data"), 890)
}

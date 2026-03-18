package fixtures

import "bytes"

func useWriteByteInvalid() {
	var buf bytes.Buffer
	buf.WriteRune('\n') // MATCH /use WriteByte instead of WriteRune for single-byte rune argument/
	buf.WriteRune('a')  // MATCH /use WriteByte instead of WriteRune for single-byte rune argument/
	buf.WriteRune('\t') // MATCH /use WriteByte instead of WriteRune for single-byte rune argument/
	buf.WriteRune('Z')  // MATCH /use WriteByte instead of WriteRune for single-byte rune argument/
	buf.WriteRune('0')  // MATCH /use WriteByte instead of WriteRune for single-byte rune argument/
}

func useWriteByteValid() {
	var buf bytes.Buffer
	// Using WriteByte directly is fine
	buf.WriteByte('\n')

	// Multi-byte rune is fine for WriteRune
	buf.WriteRune('\u00e9') // e-acute, multi-byte

	// Variable argument is fine
	r := 'a'
	buf.WriteRune(r)
}

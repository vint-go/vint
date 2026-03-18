package fixtures

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
)

func overlappingHexEncode() {
	buf := make([]byte, 100)
	hex.Encode(buf, buf[:50]) // MATCH /overlapping dst and src slices passed to an encoder/
}

func overlappingHexDecode() {
	buf := make([]byte, 100)
	hex.Decode(buf, buf[:50]) // MATCH /overlapping dst and src slices passed to an encoder/
}

func overlappingSameBuf() {
	buf := make([]byte, 100)
	hex.Encode(buf, buf) // MATCH /overlapping dst and src slices passed to an encoder/
}

func overlappingBase64Encode() {
	buf := make([]byte, 100)
	base64.StdEncoding.Encode(buf, buf[:50]) // MATCH /overlapping dst and src slices passed to an encoder/
}

func overlappingBase32Encode() {
	buf := make([]byte, 100)
	base32.StdEncoding.Encode(buf, buf[:50]) // MATCH /overlapping dst and src slices passed to an encoder/
}

func validHexEncode() {
	src := []byte("hello")
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
}

func validBase64Encode() {
	src := []byte("hello")
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
}

func validDifferentVars() {
	a := make([]byte, 100)
	b := make([]byte, 100)
	hex.Encode(a, b)
}

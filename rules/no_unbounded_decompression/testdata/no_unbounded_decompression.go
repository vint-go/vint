package fixtures

import (
	"bytes"
	"compress/bzip2"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
	"os"
)

// Invalid: io.Copy from gzip reader
func decompressGzip(src io.Reader) error {
	gz, err := gzip.NewReader(src)
	if err != nil {
		return err
	}
	defer gz.Close()

	_, err = io.Copy(os.Stdout, gz) // MATCH /use io.CopyN or io.LimitReader to prevent decompression bomb attacks/
	return err
}

// Invalid: io.Copy from zlib reader
func decompressZlib(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r) // MATCH /use io.CopyN or io.LimitReader to prevent decompression bomb attacks/
	return buf.Bytes(), err
}

// Invalid: io.Copy from flate reader
func decompressFlate(data []byte) ([]byte, error) {
	r := flate.NewReader(bytes.NewReader(data))
	defer r.Close()

	var buf bytes.Buffer
	_, err := io.Copy(&buf, r) // MATCH /use io.CopyN or io.LimitReader to prevent decompression bomb attacks/
	return buf.Bytes(), err
}

// Invalid: io.Copy from bzip2 reader
func decompressBzip2(data []byte) ([]byte, error) {
	r := bzip2.NewReader(bytes.NewReader(data))

	var buf bytes.Buffer
	_, err := io.Copy(&buf, r) // MATCH /use io.CopyN or io.LimitReader to prevent decompression bomb attacks/
	return buf.Bytes(), err
}

// Valid: io.CopyN with gzip reader (bounded)
func decompressGzipBounded(src io.Reader) error {
	gz, err := gzip.NewReader(src)
	if err != nil {
		return err
	}
	defer gz.Close()

	_, err = io.CopyN(os.Stdout, gz, 1024*1024*100)
	return err
}

// Valid: io.LimitReader wrapping gzip reader
func decompressGzipLimited(src io.Reader) ([]byte, error) {
	gz, err := gzip.NewReader(src)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	limited := io.LimitReader(gz, 1024*1024*100)
	return io.ReadAll(limited)
}

// Valid: io.Copy with a non-decompression reader
func copyNormal(src io.Reader) error {
	_, err := io.Copy(os.Stdout, src)
	return err
}

// Valid: io.Copy with os.File reader (not a decompression reader)
func copyFromFile(f *os.File) error {
	_, err := io.Copy(os.Stdout, f)
	return err
}

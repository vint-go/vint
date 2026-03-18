package pkg

import (
	"bytes"
	"encoding/binary"
)

func invalidCases() {
	var buf bytes.Buffer

	// int is not a fixed-size type
	var x int
	binary.Write(&buf, binary.LittleEndian, x) // MATCH /unsupported type int for encoding/binary, must be a fixed-size type/

	// string is not a fixed-size type
	var s string
	binary.Write(&buf, binary.LittleEndian, s) // MATCH /unsupported type string for encoding/binary, must be a fixed-size type/

	// map is not supported
	var m map[string]int
	binary.Write(&buf, binary.LittleEndian, m) // MATCH /unsupported type map[string]int for encoding/binary, must be a fixed-size type/

	// uint is not a fixed-size type
	var u uint
	binary.Write(&buf, binary.LittleEndian, u) // MATCH /unsupported type uint for encoding/binary, must be a fixed-size type/

	// binary.Read with non-fixed-size type
	var y int
	binary.Read(&buf, binary.LittleEndian, &y) // MATCH /unsupported type *int for encoding/binary, must be a fixed-size type/

	// binary.Size with non-fixed-size type
	var z int
	binary.Size(z) // MATCH /unsupported type int for encoding/binary, must be a fixed-size type/
}

func validCases() {
	var buf bytes.Buffer

	// int32 is a fixed-size type
	var x int32
	binary.Write(&buf, binary.LittleEndian, x)

	// uint16 is a fixed-size type
	var u uint16
	binary.Write(&buf, binary.LittleEndian, u)

	// float64 is a fixed-size type
	var f float64
	binary.Write(&buf, binary.LittleEndian, f)

	// bool is a fixed-size type
	var b bool
	binary.Write(&buf, binary.LittleEndian, b)

	// pointer to a fixed-size type
	var p *int32
	binary.Read(&buf, binary.LittleEndian, p)

	// slice of fixed-size types
	var sl []int32
	binary.Read(&buf, binary.LittleEndian, &sl)

	// binary.Size with fixed-size type
	var sz int64
	binary.Size(sz)
}

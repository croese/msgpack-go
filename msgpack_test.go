package msgpack

import (
	"math"
	"testing"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected []byte
	}{
		{input: nil, expected: []byte{0xc0}},
		{input: true, expected: []byte{0xc3}},
		{input: false, expected: []byte{0xc2}},
		{input: 0, expected: []byte{0}},
		{input: uint8(4), expected: []byte{0x04}},
		{input: 127, expected: []byte{0x7f}},
		{input: 128, expected: []byte{0xcc, 0x80}},
		{input: uint8(129), expected: []byte{0xcc, 0x81}},
		{input: 255, expected: []byte{0xcc, 0xff}},
		{input: 256, expected: []byte{0xcd, 0x01, 0x00}},
		{input: uint16(257), expected: []byte{0xcd, 0x01, 0x01}},
		{input: 65535, expected: []byte{0xcd, 0xff, 0xff}},
		{input: uint16(65535), expected: []byte{0xcd, 0xff, 0xff}},
		{input: 65536, expected: []byte{0xce, 0x00, 0x01, 0x00, 0x00}},
		{input: uint32(65537), expected: []byte{0xce, 0x00, 0x01, 0x00, 0x01}},
		{input: 4294967295, expected: []byte{0xce, 0xff, 0xff, 0xff, 0xff}},
		{input: uint32(4294967295), expected: []byte{0xce, 0xff, 0xff, 0xff, 0xff}},
		{input: 4294967296, expected: []byte{0xcf, 0x00, 0x00, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x00}},
		{input: uint64(4294967297), expected: []byte{0xcf, 0x00, 0x00, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x01}},
		{input: uint64(math.MaxUint64), expected: []byte{0xcf, 0xff, 0xff, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff}},
		{input: -1, expected: []byte{0xff}},
		{input: -32, expected: []byte{0xe0}},
		{input: -33, expected: []byte{0xd0, 0xdf}},
		{input: int8(-34), expected: []byte{0xd0, 0xde}},
		{input: -127, expected: []byte{0xd0, 0x81}},
		{input: int8(-127), expected: []byte{0xd0, 0x81}},
		{input: -128, expected: []byte{0xd1, 0xff, 0x80}},
		{input: int16(-129), expected: []byte{0xd1, 0xff, 0x7f}},
		{input: -2147483647, expected: []byte{0xd2, 0x80, 0x00, 0x00, 0x01}},
		{input: int32(-2147483646), expected: []byte{0xd2, 0x80, 0x00, 0x00, 0x02}},
		{input: -2147483648, expected: []byte{0xd3, 0xff, 0xff, 0xff, 0xff,
			0x80, 0x00, 0x00, 0x00}},
		{input: float32(3.14), expected: []byte{0xca, 0x40, 0x48, 0xf5, 0xc3}},
		{input: float32(-67.34), expected: []byte{0xca, 0xc2, 0x86, 0xae, 0x14}},
		{input: float64(278643.5467), expected: []byte{0xcb, 0x41, 0x11, 0x01, 0xce,
			0x2f, 0xd2, 0x1f, 0xf3}},
		{input: float64(-235.87), expected: []byte{0xcb, 0xc0, 0x6d, 0x7b, 0xd7,
			0x0a, 0x3d, 0x70, 0xa4}},
	}

	for _, test := range tests {
		b, e := Marshal(test.input)
		checkMarshalReturns(test.input, b, e, t)
		compareByteSlices(test.input, test.expected, b, t)
	}
}

func checkMarshalReturns(input interface{}, b []byte, e error, t *testing.T) {
	if e != nil {
		t.Fatalf("input %v: error should be nil", input)
	}

	if b == nil {
		t.Fatalf("input %v: returned bytes should not be nil", input)
	}
}

func compareByteSlices(input interface{}, expected []byte, actual []byte, t *testing.T) {
	for i, b := range expected {
		if actual[i] != b {
			t.Errorf("input %v: mismatched byte at index %d. expected=%x, actual=%x",
				input, i, b, actual[i])
		}
	}
}

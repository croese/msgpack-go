package msgpack

import (
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
	}

	for i, test := range tests {
		b, e := Marshal(test.input)
		checkMarshalReturns(i, b, e, t)
		compareByteSlices(i, test.expected, b, t)
	}
}

func checkMarshalReturns(testNumber int, b []byte, e error, t *testing.T) {
	if e != nil {
		t.Fatalf("test %d: error should be nil", testNumber)
	}

	if b == nil {
		t.Fatalf("test %d: returned bytes should not be nil", testNumber)
	}
}

func compareByteSlices(testNumber int, left []byte, right []byte, t *testing.T) {
	for i, b := range left {
		if right[i] != b {
			t.Errorf("test %d: mismatched byte at index %d. left=%x, right=%x",
				testNumber, i, b, right[i])
		}
	}
}

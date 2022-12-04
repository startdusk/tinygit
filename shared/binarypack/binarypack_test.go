package binarypack

import (
	"encoding/binary"
	"reflect"
	"testing"
)

func TestBinaryPack_CalcSize(t *testing.T) {
	cases := []struct {
		in   []string
		want int
		e    bool
	}{
		{[]string{}, 0, false},
		{[]string{"I", "I", "I", "4s"}, 16, false},
		{[]string{"H", "H", "I", "H", "8s", "H"}, 20, false},
		{[]string{"i", "?", "H", "f", "d", "h", "I", "5s"}, 30, false},
		{[]string{"?", "h", "H", "i", "I", "l", "L", "q", "Q", "f", "d", "1s"}, 50, false},
		// Unknown tokens
		{[]string{"a", "b", "c"}, 0, true},
	}

	for _, c := range cases {
		got, err := CalcSize(c.in)

		if err != nil && !c.e {
			t.Errorf("CalcSize(%v) raised %v", c.in, err)
		}

		if err == nil && got != c.want {
			t.Errorf("CalcSize(%v) == %d want %d", c.in, got, c.want)
		}
	}
}

func TestBinaryPack_Pack(t *testing.T) {
	cases := []struct {
		f    []string
		a    []any
		want []byte
		e    bool
	}{
		{[]string{"?", "?"}, []any{true, false}, []byte{1, 0}, false},
		{[]string{"h", "h", "h"}, []any{0, 5, -5},
			[]byte{0, 0, 5, 0, 251, 255}, false},
		{[]string{"H", "H", "H"}, []any{0, 5, 2300},
			[]byte{0, 0, 5, 0, 252, 8}, false},
		{[]string{"i", "i", "i"}, []any{0, 5, -5},
			[]byte{0, 0, 0, 0, 5, 0, 0, 0, 251, 255, 255, 255}, false},
		{[]string{"I", "I", "I"}, []any{0, 5, 2300},
			[]byte{0, 0, 0, 0, 5, 0, 0, 0, 252, 8, 0, 0}, false},
		{[]string{"f", "f", "f"}, []any{float32(0.0), float32(5.3), float32(-5.3)},
			[]byte{0, 0, 0, 0, 154, 153, 169, 64, 154, 153, 169, 192}, false},
		{[]string{"d", "d", "d"}, []any{0.0, 5.3, -5.3},
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 51, 51, 51, 51, 51, 51, 21, 64, 51, 51, 51, 51, 51, 51, 21, 192}, false},
		{[]string{"1s", "2s", "10s"}, []any{"a", "bb", "1234567890"},
			[]byte{97, 98, 98, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48}, false},
		{[]string{"I", "I", "I", "4s"}, []any{1, 2, 4, "DUMP"},
			[]byte{1, 0, 0, 0, 2, 0, 0, 0, 4, 0, 0, 0, 68, 85, 77, 80}, false},
		// Wrong format length
		{[]string{"I", "I", "I", "4s"}, []any{1, 4, "DUMP"}, nil, true},
		// Wrong format token
		{[]string{"I", "a", "I", "4s"}, []any{1, 2, 4, "DUMP"}, nil, true},
		// Wrong types
		{[]string{"?"}, []any{1.0}, nil, true},
		{[]string{"H"}, []any{int8(1)}, nil, true},
		{[]string{"I"}, []any{int32(2)}, nil, true},
		{[]string{"Q"}, []any{int64(3)}, nil, true},
		{[]string{"f"}, []any{float64(2.5)}, nil, true},
		{[]string{"d"}, []any{float32(2.5)}, nil, true},
		{[]string{"1s"}, []any{'a'}, nil, true},
	}

	for _, c := range cases {
		got, err := Pack(c.f, c.a, binary.LittleEndian)

		if err != nil && !c.e {
			t.Errorf("Pack(%v, %v) raised %v", c.f, c.a, err)
		}

		if err == nil && !reflect.DeepEqual(got, c.want) {
			t.Errorf("Pack(%v, %v) == %v want %v", c.f, c.a, got, c.want)
		}
	}
}

func TestBinaryPack_UnPack(t *testing.T) {
	cases := []struct {
		f    []string
		a    []byte
		want []any
		e    bool
	}{
		{[]string{"?", "?"}, []byte{1, 0}, []any{true, false}, false},
		{[]string{"h", "h", "h"}, []byte{0, 0, 5, 0, 251, 255},
			[]any{0, 5, -5}, false},
		{[]string{"H", "H", "H"}, []byte{0, 0, 5, 0, 252, 8},
			[]any{0, 5, 2300}, false},
		{[]string{"i", "i", "i"}, []byte{0, 0, 0, 0, 5, 0, 0, 0, 251, 255, 255, 255},
			[]any{0, 5, -5}, false},
		{[]string{"I", "I", "I"}, []byte{0, 0, 0, 0, 5, 0, 0, 0, 252, 8, 0, 0},
			[]any{0, 5, 2300}, false},
		{[]string{"f", "f", "f"},
			[]byte{0, 0, 0, 0, 154, 153, 169, 64, 154, 153, 169, 192},
			[]any{float32(0.0), float32(5.3), float32(-5.3)}, false},
		{[]string{"d", "d", "d"},
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 51, 51, 51, 51, 51, 51, 21, 64, 51, 51, 51, 51, 51, 51, 21, 192},
			[]any{0.0, 5.3, -5.3}, false},
		{[]string{"1s", "2s", "10s"},
			[]byte{97, 98, 98, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48},
			[]any{"a", "bb", "1234567890"}, false},
		{[]string{"I", "I", "I", "4s"},
			[]byte{1, 0, 0, 0, 2, 0, 0, 0, 4, 0, 0, 0, 68, 85, 77, 80},
			[]any{1, 2, 4, "DUMP"}, false},
		// Wrong format length
		{[]string{"I", "I", "I", "4s", "H"}, []byte{1, 0, 0, 0, 2, 0, 0, 0, 4, 0, 0, 0, 68, 85, 77, 80},
			nil, true},
		// Wrong format token
		{[]string{"I", "a", "I", "4s"}, []byte{1, 0, 0, 0, 2, 0, 0, 0, 4, 0, 0, 0, 68, 85, 77, 80},
			nil, true},
	}

	for _, c := range cases {
		got, err := UnPack(c.f, c.a, binary.LittleEndian)

		if err != nil && !c.e {
			t.Errorf("UnPack(%v, %v) raised %v", c.f, c.a, err)
		}

		if err == nil && !reflect.DeepEqual(got, c.want) {
			t.Errorf("UnPack(%v, %v) == %v want %v", c.f, c.a, got, c.want)
		}
	}
}

func TestBinaryPackPartialRead(t *testing.T) {
	cases := []struct {
		f    []string
		a    []byte
		i    int // Position of expected value
		want any
		e    bool
	}{
		{[]string{"I", "I", "I"}, // []any{1, 2, 4, "DUMP"} <- encoded collection has 4 values
			[]byte{1, 0, 0, 0, 2, 0, 0, 0, 4, 0, 0, 0, 68, 85, 77, 80}, 2, 4, false},
	}

	for _, c := range cases {
		got, err := UnPack(c.f, c.a, binary.LittleEndian)

		if err != nil && !c.e {
			t.Errorf("UnPack(%v, %v) raised %v", c.f, c.a, err)
		}

		if err == nil && got[c.i] != c.want {
			t.Errorf("UnPack(%v, %v) == %v want %v", c.f, c.a, got[c.i], c.want)
		}
	}
}

func TestBinaryPackUsageExample(t *testing.T) {
	// Prepare format (slice of strings)
	format := []string{"I", "?", "d", "6s"}

	// Prepare values to pack
	values := []any{4, true, 3.14, "Golang"}

	// Pack values to struct
	data, _ := Pack(format, values, binary.LittleEndian)

	// Unpack binary data to []any
	unpackedValues, _ := UnPack(format, data, binary.LittleEndian)

	if !reflect.DeepEqual(unpackedValues, values) {
		t.Errorf("Unpacked %v != original %v", unpackedValues, values)
	}

	// You can calculate size of expected binary data by format
	size, _ := CalcSize(format)

	if size != len(data) {
		t.Errorf("Size(%v) != %v", size, len(data))
	}
}

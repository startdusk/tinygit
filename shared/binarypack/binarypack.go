package binarypack

/*
	Package binarypack performs conversions between some Go values represented as byte slices.
	This can be used in handling binary data stored in files or from network connections,
	among other sources. It uses format slices of strings as compact descriptions of the layout
	of the Go structs.
	Format characters (some characters like H have been reserved for future implementation of unsigned numbers):
		? - bool, packed size 1 byte
		h, H - int, packed size 2 bytes (in future it will support pack/unpack of int8, uint8 values)
		i, I, l, L - int, packed size 4 bytes (in future it will support pack/unpack of int16, uint16, int32, uint32 values)
		q, Q - int, packed size 8 bytes (in future it will support pack/unpack of int64, uint64 values)
		f - float32, packed size 4 bytes
		d - float64, packed size 8 bytes
		Ns - string, packed size N bytes, N is a number of runes to pack/unpack
*/

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Return a byte slice containing the values of msg slice packed according to the given format.
// The items of msg slice must match the values required by the format exactly.
func Pack(format []string, msg []any, order binary.ByteOrder) ([]byte, error) {
	if len(format) > len(msg) {
		return nil, errors.New("format is longer than values to pack")
	}

	var res []byte

	for i, f := range format {
		switch f {
		case "?":
			castedValue, ok := msg[i].(bool)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (bool)")
			}
			res = append(res, boolToBytes(castedValue)...)
		case "h", "H":
			castedValue, ok := msg[i].(int)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (int, 2 bytes)")
			}
			res = append(res, intToBytes(castedValue, 2)...)
		case "i", "I", "l", "L":
			castedValue, ok := msg[i].(int)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (int, 4 bytes)")
			}
			res = append(res, intToBytes(castedValue, 4)...)
		case "q", "Q":
			castedValue, ok := msg[i].(int)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (int, 8 bytes)")
			}
			res = append(res, intToBytes(castedValue, 8)...)
		case "f":
			castedValue, ok := msg[i].(float32)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (float32)")
			}
			res = append(res, float32ToBytes(castedValue, 4, order)...)
		case "d":
			castedValue, ok := msg[i].(float64)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (float64)")
			}
			res = append(res, float64ToBytes(castedValue, 8, order)...)
		default:
			if strings.Contains(f, "s") {
				castedValue, ok := msg[i].(string)
				if !ok {
					return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (string)")
				}
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				res = append(res, []byte(fmt.Sprintf("%s%s",
					castedValue, strings.Repeat("\x00", n-len(castedValue))))...)
			} else {
				return nil, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}

	return res, nil
}

// Unpack the byte slice (presumably packed by Pack(format, msg)) according to the given format.
// The result is a []any slice even if it contains exactly one item.
// The byte slice must contain not less the amount of data required by the format
// (len(msg) must more or equal CalcSize(format)).
func UnPack(format []string, msg []byte, order binary.ByteOrder) ([]any, error) {
	expectedSize, err := CalcSize(format)

	if err != nil {
		return nil, err
	}

	if expectedSize > len(msg) {
		return nil, errors.New("expected size is bigger than actual size of message")
	}

	var res []any

	for _, f := range format {
		switch f {
		case "?":
			res = append(res, bytesToBool(msg[:1], order))
			msg = msg[1:]
		case "h", "H":
			res = append(res, bytesToInt(msg[:2], order))
			msg = msg[2:]
		case "i", "I", "l", "L":
			res = append(res, bytesToInt(msg[:4], order))
			msg = msg[4:]
		case "q", "Q":
			res = append(res, bytesToInt(msg[:8], order))
			msg = msg[8:]
		case "f":
			res = append(res, bytesToFloat32(msg[:4], order))
			msg = msg[4:]
		case "d":
			res = append(res, bytesToFloat64(msg[:8], order))
			msg = msg[8:]
		default:
			if strings.Contains(f, "s") {
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				res = append(res, string(msg[:n]))
				msg = msg[n:]
			} else {
				return nil, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}

	return res, nil
}

// Return the size of the struct (and hence of the byte slice) corresponding to the given format.
func CalcSize(format []string) (int, error) {
	var size int

	for _, f := range format {
		switch f {
		case "?":
			size = size + 1
		case "h", "H":
			size = size + 2
		case "i", "I", "l", "L", "f":
			size = size + 4
		case "q", "Q", "d":
			size = size + 8
		default:
			if strings.Contains(f, "s") {
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				size = size + n
			} else {
				return 0, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}

	return size, nil
}

func boolToBytes(x bool) []byte {
	if x {
		return intToBytes(1, 1)
	}
	return intToBytes(0, 1)
}

func bytesToBool(b []byte, order binary.ByteOrder) bool {
	return bytesToInt(b, order) > 0
}

func intToBytes(n int, size int) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, int64(n))
	return buf.Bytes()[0:size]
}

func bytesToInt(b []byte, order binary.ByteOrder) int {
	buf := bytes.NewBuffer(b)

	switch len(b) {
	case 1:
		var x int8
		binary.Read(buf, order, &x)
		return int(x)
	case 2:
		var x int16
		binary.Read(buf, order, &x)
		return int(x)
	case 4:
		var x int32
		binary.Read(buf, order, &x)
		return int(x)
	default:
		var x int64
		binary.Read(buf, order, &x)
		return int(x)
	}
}

func float32ToBytes(n float32, size int, order binary.ByteOrder) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, order, n)
	return buf.Bytes()[0:size]
}

func bytesToFloat32(b []byte, order binary.ByteOrder) float32 {
	var x float32
	buf := bytes.NewBuffer(b)
	binary.Read(buf, order, &x)
	return x
}

func float64ToBytes(n float64, size int, order binary.ByteOrder) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, order, n)
	return buf.Bytes()[0:size]
}

func bytesToFloat64(b []byte, order binary.ByteOrder) float64 {
	var x float64
	buf := bytes.NewBuffer(b)
	binary.Read(buf, order, &x)
	return x
}

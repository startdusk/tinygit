package tinygit

import (
	"reflect"
	"testing"
)

func TestIndexLifecycle(t *testing.T) {
	indexes := []Index{
		{
			CTimeS: 1670330587083,
			CTimeN: 0,
			MTimeS: 1670330580274,
			MTimeN: 0,
			Dev:    56,
			INO:    44,
			Mode:   20,
			UID:    134,
			GID:    12,
			Size:   1 << 20,
			Sha1:   "20stringlen111111919",
			Flags:  0,
			Path:   "/path/to/test.txt",
		},
	}

	err := WriteIndex(indexes)
	if err != nil {
		t.Fatalf("write index: %+v", err)
	}

	readed, err := ReadIndex()
	if err != nil {
		t.Fatalf("read index: %+v", err)
	}

	if len(readed) != len(indexes) {
		t.Fatalf("write count %d, but read count %d", len(indexes), len(readed))
	}

	if !reflect.DeepEqual(readed[0], indexes[0]) {
		t.Fatalf("write index %+v, but read index %+v", indexes[0], readed[0])
	}
}

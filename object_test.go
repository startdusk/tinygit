package tinygit

import (
	"reflect"
	"testing"
)

func TestHashObject(t *testing.T) {
	sha1, err := HashObject(HashParam{
		Data:      []byte("Hello world"),
		ObjType:   Blob,
		WriteFile: true,
	})
	if err != nil || sha1 == "" {
		t.Fatal(err)
	}
}

func TestZlibFunction(t *testing.T) {
	content := []byte("Hello World!")
	compressed, err := zlibCompress(content)
	t.Logf("content=%s, compress=%+v", content, compressed)
	if err != nil {
		t.Fatalf("zlib compress: %+v", err)
	}

	decompress, err := zlibDecompress(compressed)
	if err != nil {
		t.Fatalf("zlib decompress: %+v", err)
	}
	if !reflect.DeepEqual(content, decompress) {
		t.Fatalf("content=%s, decompress=%s", content, decompress)
	}
}

func FuzzZlibCompress(f *testing.F) {
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, a []byte) {
		_, err := zlibCompress(a)
		if err != nil {
			t.Fatalf("zlibCompress: %+v", err)
		}
	})
}

package tinygit

import (
	"reflect"
	"testing"
)

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

func FuzzZlibFunction(f *testing.F) {
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, content []byte) {
		compressed, err := zlibCompress(content)
		if err != nil {
			t.Fatalf("zlibCompress: %+v", err)
		}
		decompress, err := zlibDecompress(compressed)
		if err != nil {
			t.Fatalf("zlib decompress: %+v", err)
		}
		if !reflect.DeepEqual(content, decompress) {
			t.Fatalf("content=%s, decompress=%s", content, decompress)
		}
	})
}

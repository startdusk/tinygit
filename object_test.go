package tinygit

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestHashObject(t *testing.T) {
	content := "Hello World"
	data := []byte(content)
	sha1, objPath, err := HashObject(HashParam{
		Data:      data,
		ObjType:   Blob,
		WriteFile: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	if sha1 == "" {
		t.Fatal("sha1 empty")
	}

	if objPath == "" {
		t.Fatal("obj path empty")
	}

	objPathList := strings.Split(objPath, "/")
	if len(objPathList) != 4 {
		t.Fatalf("expected path split 4 pair, but got %d, source path %s", len(objPathList), objPath)
	}
	if objPathList[0] != RepoRootPath {
		t.Fatalf("expected root path %s, but got %s", RepoRootPath, objPathList[0])
	}
	if objPathList[1] != ObjectsFolder {
		t.Fatalf("expected object folder %s, but got %s", ObjectsFolder, objPathList[1])
	}
	if objPathList[2]+objPathList[3] != sha1 {
		t.Fatalf("expected sha1 %s, but got %s", sha1, objPathList[2]+objPathList[3])
	}

	compressed, err := os.ReadFile(objPath)
	if err != nil {
		t.Fatalf("cannot open the file %s: %+v", objPath, err)
	}

	decompressed, err := zlibDecompress(compressed)
	if err != nil {
		t.Fatalf("cannot decompress file: %+v", err)
	}
	expectDecompressed := []byte{98, 108, 111, 98, 32, 49, 49, 0, 72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100}
	if !reflect.DeepEqual(expectDecompressed, decompressed) {
		t.Fatalf("expected decompressed %v, but got %v", expectDecompressed, decompressed)
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

func TestMain(m *testing.M) {
	defer func() {
		os.RemoveAll(RepoRootPath)
	}()
	m.Run()
}

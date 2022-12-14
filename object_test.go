package tinygit

import (
	"os"
	"path/filepath"
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

	objPath = filepath.ToSlash(objPath)
	objPathList := strings.Split(objPath, "/")
	if len(objPathList) != 4 {
		t.Fatalf("expected path split 4 pair, but got %d, source path %s",
			len(objPathList), objPath)
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

func TestObjectLifecycle(t *testing.T) {

	cases := []struct {
		name      string
		param     HashParam
		writeFile bool
	}{
		{
			name: "hash_object_write_file",
			param: HashParam{
				ObjType:   Blob,
				Data:      []byte("hello world"),
				WriteFile: true,
			},
			writeFile: true,
		},
		{
			name: "hash_object_no_write_file",
			param: HashParam{
				ObjType:   Blob,
				Data:      []byte("hello world111"),
				WriteFile: false,
			},
			writeFile: false,
		},

		// TODO: need more test cases
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// hash object
			sha1, objFile, err := HashObject(c.param)
			if err != nil {
				t.Fatalf("hash object: %+v", err)
			}

			if c.writeFile {
				// find object
				gotObjFile, err := FindObject(sha1)
				if err != nil {
					t.Fatalf("find object: %+v", err)
				}

				if gotObjFile != objFile {
					t.Fatalf("cannot find the object file [%s], but got [%s]", objFile, gotObjFile)
				}

				// read object
				obj, err := ReadObject(sha1)
				if err != nil {
					t.Fatalf("read object: %+v", err)
				}
				if obj.Size != len(c.param.Data) {
					t.Fatalf("expected data size %d, but got %d", len(c.param.Data), obj.Size)
				}

				if obj.Type != c.param.ObjType {
					t.Fatalf("expected data type %s, but got %s", c.param.ObjType, obj.Type)
				}

				if !reflect.DeepEqual(obj.Data, c.param.Data) {
					t.Fatalf("expected data %v, but got %v", c.param.Data, obj.Data)
				}
			} else {
				if objFile != "" {
					t.Fatalf("expected no write file, but got %s", objFile)
				}
			}
		})
	}
}

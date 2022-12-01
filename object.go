package tinygit

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	RepoRootPath  = ".tinygit"
	ObjectsFolder = "objects"
)

type ObjType = string

// There are three types of objects in the Git model: blobs (ordinary files),
// commits, and trees (these represent the state of a single directory).
const (
	Blob   ObjType = "blob"
	Commit ObjType = "commit"
	Tree   ObjType = "tree"
)

type HashParam struct {
	Data      []byte
	ObjType   ObjType
	WriteFile bool
}

// HashObject compute hash of object data of given type and write to object store if
// "WriteFile" is True. Return SHA-1 object hash as hex string.
func HashObject(param HashParam) (string, string, error) {
	fullData := genFullData(param.ObjType, param.Data)
	sha1 := sha1Hash(fullData)

	var objPath string
	if param.WriteFile {
		path := genObjectPath(sha1[:2])
		// check if file or directory exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// path does not exist
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return sha1, objPath, fmt.Errorf("create path: %w", err)
			}
		}
		compressed, err := zlibCompress(fullData)
		if err != nil {
			return sha1, objPath, fmt.Errorf("zlib compress: %w", err)
		}
		objPath = genObjectFile(path, sha1[2:])
		if err := os.WriteFile(objPath, compressed, os.ModePerm); err != nil {
			return sha1, objPath, fmt.Errorf("write file: %w", err)
		}
	}
	return sha1, objPath, nil
}

// FindObject find object with given SHA-1 prefix and return path to object in object
// store, or raise ValueError if there are no objects or multiple objects
// with this prefix.
func FindObject(sha1Prefix string) error {
	return nil
}

func sha1Hash(data []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(data))
}

func zlibCompress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer w.Close()
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	w.Flush()
	return b.Bytes(), nil
}

func zlibDecompress(compressed []byte) ([]byte, error) {
	b := bytes.NewReader(compressed)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.Bytes(), nil
}

func genFullData(typ ObjType, data []byte) []byte {
	header := []byte(fmt.Sprintf("%s %d", typ, len(data)))
	var fullData []byte
	fullData = append(fullData, header...)
	fullData = append(fullData, ' ')
	fullData = append(fullData, data...)
	return fullData
}

func genObjectFile(path, filename string) string {
	return filepath.Join(path, filename)
}

func genObjectPath(filename string) string {
	return filepath.Join(RepoRootPath, ObjectsFolder, filename)
}

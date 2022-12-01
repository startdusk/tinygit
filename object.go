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
func HashObject(param HashParam) (string, error) {
	header := []byte(fmt.Sprintf("%s %d", param.ObjType, len(param.Data)))
	var fullData []byte
	fullData = append(fullData, header...)
	fullData = append(fullData, '\x00')
	fullData = append(fullData, param.Data...)
	sha1 := sha1Hash(fullData)
	if param.WriteFile {
		path := genObjectPath(sha1[:2])
		// check if file or directory exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// path does not exist
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return sha1, fmt.Errorf("create path: %w", err)
			}
		}
		compress, err := zlibCompress(fullData)
		if err != nil {
			return sha1, fmt.Errorf("zlib compress: %w", err)
		}
		if err := os.WriteFile(genObjectFile(path, sha1[2:]), compress, os.ModePerm); err != nil {
			return sha1, fmt.Errorf("write file: %w", err)
		}
	}
	return sha1, nil
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

func genObjectFile(path, filename string) string {
	return filepath.Join(path, filename)
}

func genObjectPath(filename string) string {
	return filepath.Join(RepoRootPath, ObjectsFolder, filename)
}

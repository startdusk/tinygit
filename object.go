package tinygit

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	RepoRootPath  = ".tinygit"
	ObjectsFolder = "objects"
)

type ObjType = string

// Object represents a tinygit object.
type Object struct {
	Type ObjType
	Size int
	Data []byte
}

// NewObject creates a tinygit object.
func NewObject(typ ObjType, data []byte) Object {
	return Object{
		Type: typ,
		Size: len(data),
		Data: data,
	}
}

// FormatData structs data
func (o Object) FormatData() []byte {
	header := []byte(fmt.Sprintf("%s %d", o.Type, o.Size))
	var fullData []byte
	fullData = append(fullData, header...)
	fullData = append(fullData, '\x00')
	fullData = append(fullData, o.Data...)
	return fullData
}

func (o Object) String() string {
	return string(o.FormatData())
}

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
	fullData := NewObject(param.ObjType, param.Data).FormatData()
	sha1 := sha1Hash(fullData)

	var objFile string
	if param.WriteFile {
		path := genObjectPath(sha1[:2])
		// check if file or directory exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// path does not exist
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return sha1, objFile, fmt.Errorf("create path: %w", err)
			}
		}
		compressed, err := zlibCompress(fullData)
		if err != nil {
			return sha1, objFile, fmt.Errorf("zlib compress: %w", err)
		}
		objFile = genObjectFile(path, sha1[2:])
		if err := os.WriteFile(objFile, compressed, os.ModePerm); err != nil {
			return sha1, objFile, fmt.Errorf("write file: %w", err)
		}
	}
	return sha1, objFile, nil
}

// FindObject find object with given SHA-1 prefix and return path to object in object
// store, or raise ValueError if there are no objects or multiple objects
// with this prefix.
func FindObject(sha1Prefix string) (string, error) {
	if len(sha1Prefix) < 2 {
		return "", errors.New("hash prefix must be 2 or more characters")
	}

	objPath := genObjectPath(sha1Prefix[:2])
	objFile := genObjectFile(objPath, sha1Prefix[2:])
	_, err := os.Stat(objFile)
	if err != nil {
		return "", fmt.Errorf("find object: %w", err)
	}

	return objFile, nil
}

// Read object with given SHA-1 prefix
func ReadObject(sha1Prefix string) (Object, error) {
	// 1.find object
	path, err := FindObject(sha1Prefix)
	if err != nil {
		return Object{}, err
	}

	// 2.open object
	compressed, err := os.ReadFile(path)
	if err != nil {
		return Object{}, fmt.Errorf("read file: %w", err)
	}

	// 3.parse object
	return parseObject(compressed)
}

func parseObject(compressed []byte) (Object, error) {
	fullData, err := zlibDecompress(compressed)
	if err != nil {
		return Object{}, err
	}
	nullIndex := bytes.Index(fullData, []byte("\x00"))
	header := fullData[:nullIndex]
	headerSplit := strings.Split(string(header), " ")
	if len(headerSplit) != 2 {
		return Object{}, fmt.Errorf("object header should have 2 pair but got %d", len(headerSplit))
	}
	objType := ObjType(headerSplit[0])
	size, err := strconv.Atoi(headerSplit[1])
	if err != nil {
		return Object{}, fmt.Errorf("size should be a number, but got %s: %w", headerSplit[1], err)
	}
	data := fullData[nullIndex:]
	if len(data) != size {
		return Object{}, fmt.Errorf("expected size %d, got %d bytes", size, len(data))
	}
	return Object{Type: objType, Data: data}, nil
}

func sha1Hash(data []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(data))
}

func genObjectFile(path, filename string) string {
	return filepath.Join(path, filename)
}

func genObjectPath(filename string) string {
	return filepath.Join(RepoRootPath, ObjectsFolder, filename)
}

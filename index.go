package tinygit

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"

	"github.com/startdusk/tinygit/shared/binarypack"
)

// Index represents a index struct.
type Index struct {
	CTimeS int64
	CTimeN int64
	MTimeS int64
	MTimeN int64
	Dev    int32
	INO    uint64
	Mode   uint16
	UID    uint32
	GID    uint32
	Size   int64
	Sha1   string
	Flags  uint32
	Path   string
}

// Head pack the head message.
func (i Index) Head() ([]byte, error) {
	values := []any{
		i.CTimeS,
		i.CTimeN,
		i.MTimeS,
		i.MTimeN,
		i.Dev,
		i.INO,
		i.Mode,
		i.UID,
		i.GID,
		i.Size,
		i.Sha1,
		i.Flags,
	}
	return binarypack.Pack(headFormat, values, binary.BigEndian)
}

// Indexes represents a index slice for sort.
type Indexes []Index

func (idxs Indexes) Len() int           { return len(idxs) }
func (idxs Indexes) Swap(i, j int)      { idxs[i], idxs[j] = idxs[j], idxs[i] }
func (idxs Indexes) Less(i, j int) bool { return idxs[i].Path < idxs[j].Path }

// Sort sort the index.
func (idxs Indexes) Sort() []Index {
	sort.Sort(idxs)
	return idxs
}

const (
	indexSignature = "DIRC"
	indexVersion   = int32(1)

	headSize     = 108
	headerSize   = 12
	checkSumSize = 40
)

var indexFile = filepath.Join(RepoRootPath, "index")

var headFormat = []string{"Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "Q", "20s", "Q"}

var headerFormat = []string{"4s", "L", "L"}

// ReadIndex read tinygit index file and return list of Index objects.
func ReadIndex() (Indexes, error) {
	indexes := make([]Index, 0)
	data, err := os.ReadFile(indexFile)
	if err != nil {
		switch {
		case errors.Is(err, fs.ErrNotExist):
			return indexes, nil
		default:
			return nil, fmt.Errorf("read index file: %w", err)
		}
	}

	digest := sha1Hash(data[:len(data)-checkSumSize])
	if !reflect.DeepEqual([]byte(digest), data[len(data)-checkSumSize:]) {
		return nil, errors.New("invalid index checksum")
	}

	header, err := binarypack.UnPack(headerFormat, data[:headerSize], binary.LittleEndian)
	if err != nil {
		return nil, err
	}

	if len(header) != 3 {
		return nil, fmt.Errorf("invalid header: %v", header)
	}

	signature, ok := header[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid signature type")
	}
	if signature != indexSignature {
		return nil, fmt.Errorf("invalid signature: `%s`", signature)
	}

	version, ok := header[1].(int)
	if !ok {
		return nil, fmt.Errorf("invalid version type")
	}
	if int32(version) != indexVersion {
		return nil, fmt.Errorf("unknown index version %d", version)
	}

	indexNum, ok := header[2].(int)
	if !ok {
		return nil, fmt.Errorf("invalid index object lenght type")
	}

	indexesBytes := data[headerSize : len(data)-checkSumSize]

	var i int
	for i+headSize < len(indexesBytes) {
		fieldsEnd := i + headSize
		fields, err := binarypack.UnPack(headFormat, indexesBytes[i:fieldsEnd], binary.LittleEndian)
		if err != nil {
			return nil, err
		}
		pathEnd := bytes.IndexByte(indexesBytes[fieldsEnd:], '\x00')
		path := indexesBytes[fieldsEnd : fieldsEnd+pathEnd]
		index, err := parseFieldsToIndex(fields)
		if err != nil {
			return nil, err
		}
		index.Path = string(path)
		indexes = append(indexes, index)
		indexLen := ((headSize + len(path) + 8) / 8) * 8
		i += indexLen
	}
	if len(indexes) != int(indexNum) {
		return nil, fmt.Errorf("invalid index num")
	}
	return indexes, nil
}

// WriteIndex write list of Index objects to tinygit index file.
func WriteIndex(indexes []Index) error {
	var packeds [][]byte
	for _, index := range indexes {
		packed, err := makePacked(index)
		if err != nil {
			return err
		}
		packeds = append(packeds, packed)
	}
	header, err := binarypack.Pack(headerFormat, []any{indexSignature, indexVersion,
		int32(len(indexes))}, binary.BigEndian)
	if err != nil {
		return err
	}

	allData := []byte(header)
	for _, packed := range packeds {
		allData = append(allData, packed...)
	}
	digest := sha1Hash(allData)
	allData = append(allData, []byte(digest)...)
	return os.WriteFile(indexFile, allData, os.ModePerm)
}

func makePacked(index Index) ([]byte, error) {
	head, err := index.Head()
	if err != nil {
		return nil, err
	}
	path := []byte(index.Path)
	length := ((headSize + len(path) + 8) / 8) * 8
	var packed []byte
	packed = append(packed, head...)
	packed = append(packed, path...)
	repeat := (length - headSize - len(path))
	var padding []byte
	if repeat > 0 {
		padding = bytes.Repeat([]byte("\x00"), repeat)
	}

	packed = append(packed, padding...)
	return packed, nil
}

func parseFieldsToIndex(data []any) (Index, error) {
	fmt.Printf("%v", data)
	index := Index{}
	if len(data) != 12 {
		return index, fmt.Errorf("invalid index params")
	}
	cTimeS, ok := data[0].(int)
	if !ok {
		return index, fmt.Errorf("invalid ctime_s type")
	}
	index.CTimeS = int64(cTimeS)

	cTimeN, ok := data[1].(int)
	if !ok {
		return index, fmt.Errorf("invalid ctime_n type")
	}
	index.CTimeN = int64(cTimeN)

	mTimeS, ok := data[2].(int)
	if !ok {
		return index, fmt.Errorf("invalid mtime_s type")
	}
	index.MTimeS = int64(mTimeS)

	mTimeN, ok := data[3].(int)
	if !ok {
		return index, fmt.Errorf("invalid mtime_n type")
	}
	index.MTimeN = int64(mTimeN)

	dev, ok := data[4].(int)
	if !ok {
		return index, fmt.Errorf("invalid dev type")
	}
	index.Dev = int32(dev)

	ino, ok := data[5].(int)
	if !ok {
		return index, fmt.Errorf("invalid ino type")
	}
	index.INO = uint64(ino)

	mode, ok := data[6].(int)
	if !ok {
		return index, fmt.Errorf("invalid mode type")
	}
	index.Mode = uint16(mode)

	uid, ok := data[7].(int)
	if !ok {
		return index, fmt.Errorf("invalid uid type")
	}
	index.UID = uint32(uid)

	gid, ok := data[8].(int)
	if !ok {
		return index, fmt.Errorf("invalid gid type")
	}
	index.GID = uint32(gid)

	size, ok := data[9].(int)
	if !ok {
		return index, fmt.Errorf("invalid size type")
	}
	index.Size = int64(size)

	sha1, ok := data[10].(string)
	if !ok {
		return index, fmt.Errorf("invalid sha1 type")
	}
	index.Sha1 = sha1

	flags, ok := data[11].(int)
	if !ok {
		return index, fmt.Errorf("invalid flags type")
	}
	index.Flags = uint32(flags)

	return index, nil
}

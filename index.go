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

// Index represens a index struct.
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

// Indexes represens a index slice for sort.
type Indexes []Index

func (idxs Indexes) Len() int           { return len(idxs) }
func (idxs Indexes) Swap(i, j int)      { idxs[i], idxs[j] = idxs[j], idxs[i] }
func (idxs Indexes) Less(i, j int) bool { return idxs[i].Path < idxs[j].Path }

func (idxs Indexes) Sort() []Index {
	sort.Sort(idxs)
	return idxs
}

const indexSignature = "DIRC"

const indexVersion = 1

var indexFile = filepath.Join(RepoRootPath, "index")
var headFormat = []string{"L", "L", "L", "L", "L", "L", "L", "L", "20s", "H"}
var headerFormat = []string{"4s", "L", "L"}

// Read tinygit index file and return list of Index objects.
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

	digest := sha1Hash(data[:len(data)-20])
	if !reflect.DeepEqual(digest, data[len(data)-20:]) {
		return nil, errors.New("invalid index checksum")
	}
	header, err := binarypack.UnPack([]string{"4s", "L", "L"}, data[:12], binary.BigEndian)
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
	if version != indexVersion {
		return nil, fmt.Errorf("unknown index version %d", version)
	}

	indexNum, ok := header[2].(int)
	if !ok {
		return nil, fmt.Errorf("invalid index object lenght type")
	}

	indexesBytes := data[12 : len(data)-20]

	var i int
	for i+62 < len(indexesBytes) {
		fieldsEnd := i + 62
		fields, err := binarypack.UnPack(headFormat, indexesBytes[i:fieldsEnd], binary.BigEndian)
		if err != nil {
			return nil, err
		}
		pathEnd := bytes.IndexByte(data[fieldsEnd:], '\x00')
		path := data[fieldsEnd:pathEnd]
		index, err := parseFieldsToIndex(fields)
		if err != nil {
			return nil, err
		}
		indexes = append(indexes, index)
		indexLen := ((62 + len(path) + 8) / 8) * 8
		i += indexLen
	}
	if len(indexes) != indexNum {
		return nil, fmt.Errorf("invalid index num")
	}
	return indexes, nil
}

// Write list of Index objects to tinygit index file.
func WriteIndex(indexes []Index) error {
	var packeds [][]byte
	for _, index := range indexes {
		packed, err := makePacked(index)
		if err != nil {
			return err
		}
		packeds = append(packeds, packed)
	}
	header, err := binarypack.Pack(headerFormat, []any{indexSignature, indexVersion, len(indexes)}, binary.BigEndian)
	if err != nil {
		return err
	}
	var allData []byte
	allData = append(allData, []byte(header)...)
	for _, packed := range packeds {
		allData = append(allData, packed...)
	}
	digest := sha1Hash(allData)
	allData = append(allData, []byte(digest)...)
	return os.WriteFile(indexFile, allData, os.ModePerm)
}

func makePacked(index Index) ([]byte, error) {
	values := []any{
		index.CTimeS,
		index.CTimeN,
		index.MTimeS,
		index.MTimeN,
		index.Dev,
		index.INO,
		index.Mode,
		index.UID,
		index.GID,
		index.Size,
		index.Sha1,
		index.Flags,
	}
	head, err := binarypack.Pack(headFormat, values, binary.BigEndian)
	if err != nil {
		return nil, err
	}

	path := []byte(index.Path)
	length := ((62 + len(path) + 8) / 8) * 8
	var packed []byte
	packed = append(packed, []byte(head)...)
	packed = append(packed, path...)
	padding := bytes.Repeat([]byte("\x00"), (length - 68 - len(path)))
	packed = append(packed, padding...)
	return packed, nil
}

func parseFieldsToIndex(data []any) (Index, error) {
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

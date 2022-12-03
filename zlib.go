package tinygit

import (
	"bytes"
	"compress/zlib"
	"io"
)

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

package tinygit

import "testing"

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

package tinygit

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/startdusk/tinygit/shared/filestat"
)

// Create directory for repo and initialize .tinygit directory.
func Initail(repo string) error {
	if err := os.MkdirAll(repo, os.ModePerm); err != nil {
		return err
	}
	tinygitPath := filepath.Join(repo, ".tinygit")
	if err := os.MkdirAll(tinygitPath, os.ModePerm); err != nil {
		return err
	}
	for _, name := range [3]string{"objects", "refs", "refs/heads"} {
		if err := os.MkdirAll(filepath.Join(tinygitPath, name), os.ModePerm); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(tinygitPath, "HEAD"), []byte("ref: refs/heads/master"), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func PrintHelp() {
	const help = `useage: tinygit [-v | --version] [-h | --help]
	        <command> [<args>]
These are common TinyGit commands used in various situations:

start a working area (see also: tinygit help tutorial)
   init      Create an empty TinyGit repository or reinitialize an existing one

work on the current change (see also: tinygit help everyday)
   add       Add file contents to the index
   mv        Move or rename a file, a directory, or a symlink
   rm        Remove files from the working tree and from the index
	`
	fmt.Println(help)
}

func Add(path string) error {
	// 1.read files recursively
	var paths []string
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path, info.Size())
		return nil
	})
	if err != nil {
		return err
	}
	// 2.read index all entries
	indexes, err := ReadIndex()
	if err != nil {
		return err
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sha1, _, err := HashObject(HashParam{
			Data:      data,
			ObjType:   Blob,
			WriteFile: true,
		})
		if err != nil {
			return err
		}
		st, err := filestat.Stat(path)
		if err != nil {
			return err
		}
		index := Index{
			CTimeS: st.CTimeS,
			CTimeN: 0,
			MTimeS: st.MTimeS,
			MTimeN: 0,
			Dev:    st.Dev,
			INO:    st.INO,
			Mode:   st.Mode,
			UID:    st.UID,
			GID:    st.GID,
			Size:   st.Size,
			Sha1:   sha1,
			Flags:  st.Flags,
			Path:   path,
		}
		indexes = append(indexes, index)
	}
	indexes = indexes.Sort()
	return WriteIndex(indexes)
}

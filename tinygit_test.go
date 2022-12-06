package tinygit

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// first create the repo
	Initail(".")
	// last remove the repo
	defer func() {
		os.RemoveAll(RepoRootPath)
	}()
	m.Run()
}

//go:build darwin

package filestat

import (
	"testing"
)

func TestFileStat(t *testing.T) {
	// TODO: Strengthen this test case.
	st, err := Stat("./filestat_darwin_test.go")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", st)
}

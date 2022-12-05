//go:build windows

package filestat

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// FileStat indicates the file status of the specified system.
type FileStat struct {
	CreateTime int64
	ModifyTime int64
	Dev        int32
	INO        uint64
	Mode       uint16
	UID        uint32
	GID        uint32
	Size       int64
	Flags      uint32
}

// Stat query file state information.
func Stat(path string) (FileStat, error) {
	filestat := FileStat{}
	fileinfo, err := os.Stat(path)
	if err != nil {
		return filestat, err
	}
	stat, ok := fileinfo.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return filestat, fmt.Errorf("invalid system stat")
	}

	filestat.CreateTime = time.Unix(0, stat.CreationTime.Nanoseconds()).Unix()
	filestat.ModifyTime = time.Unix(0, stat.LastWriteTime.Nanoseconds()).Unix()
	return filestat, nil
}

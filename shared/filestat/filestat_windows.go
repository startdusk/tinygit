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
	CTimeS int64
	MTimeS int64
	Dev    int32
	INO    uint64
	Mode   uint16
	UID    uint32
	GID    uint32
	Size   int64
	Flags  uint32
}

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
	// filestat.Dev = stat.Dev
	// filestat.INO = stat.Ino
	// filestat.Mode = stat.Mode
	// filestat.UID = stat.Uid
	// filestat.GID = stat.Gid
	// filestat.Size = stat.Size
	// filestat.Flags = stat.Flags
	filestat.CTimeS = int64(time.Since(time.Unix(0, stat.CreationTime.Nanoseconds())))
	filestat.MTimeS = int64(time.Since(time.Unix(0, stat.LastWriteTime.Nanoseconds())))
	return filestat, nil
}
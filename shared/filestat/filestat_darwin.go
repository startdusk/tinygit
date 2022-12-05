//go:build darwin

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
	stat, ok := fileinfo.Sys().(*syscall.Stat_t)
	if !ok {
		return filestat, fmt.Errorf("invalid system stat")
	}
	filestat.Dev = stat.Dev
	filestat.INO = stat.Ino
	filestat.Mode = stat.Mode
	filestat.UID = stat.Uid
	filestat.GID = stat.Gid
	filestat.Size = stat.Size
	filestat.Flags = stat.Flags
	filestat.CreateTime = time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec).Unix()
	filestat.ModifyTime = time.Unix(stat.Mtimespec.Sec, stat.Mtimespec.Nsec).Unix()
	return filestat, nil
}

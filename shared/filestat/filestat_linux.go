//go:build linux

package filestat

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// type Stat_t struct {
// 	Dev       uint64
// 	Ino       uint64
// 	Nlink     uint64
// 	Mode      uint32
// 	Uid       uint32
// 	Gid       uint32
// 	X__pad0   int32
// 	Rdev      uint64
// 	Size      int64
// 	Blksize   int64
// 	Blocks    int64
// 	Atim      Timespec
// 	Mtim      Timespec
// 	Ctim      Timespec
// 	X__unused [3]int64
// }

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
	filestat.Dev = int32(stat.Dev)
	filestat.INO = stat.Ino
	filestat.Mode = uint16(stat.Mode)
	filestat.UID = stat.Uid
	filestat.GID = stat.Gid
	filestat.Size = stat.Size
	filestat.CreateTime = time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec).Unix()
	filestat.ModifyTime = time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec).Unix()
	return filestat, nil
}

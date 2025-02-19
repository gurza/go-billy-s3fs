package s3fs

import (
	"os"
	"time"
)

// fileStat is the implementation of FileInfo returned by Stat and Lstat.
type fileStat struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	// sys     syscall.Stat_t
}

func newFileInfo(name string, size int64, modTime time.Time) os.FileInfo {
	return &fileStat{
		name:    name,
		size:    size,
		modTime: modTime,
		mode:    0644,
	}
}

func newDirInfo(name string) os.FileInfo {
	return &fileStat{
		name: name,
		mode: os.ModeDir | 0755,
	}
}

func (fs *fileStat) Name() string       { return fs.name }
func (fs *fileStat) IsDir() bool        { return fs.mode.IsDir() }
func (fs *fileStat) Size() int64        { return fs.size }
func (fs *fileStat) Mode() os.FileMode  { return fs.mode }
func (fs *fileStat) ModTime() time.Time { return fs.modTime }
func (fs *fileStat) Sys() any           { return nil }

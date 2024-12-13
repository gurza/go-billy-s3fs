package s3fs

import (
	"os"
	"time"
)

type fileInfo struct {
	name    string
	size    int64
	modTime time.Time
	mode    os.FileMode
}

func newFileInfo(name string, size int64, modTime time.Time) os.FileInfo {
	return &fileInfo{
		name:    name,
		size:    size,
		modTime: modTime,
		mode:    0644,
	}
}

func newDirInfo(name string) os.FileInfo {
	return &fileInfo{
		name: name,
		mode: os.ModeDir | 0755,
	}
}

func (f *fileInfo) Name() string       { return f.name }
func (f *fileInfo) Size() int64        { return f.size }
func (f *fileInfo) Mode() os.FileMode  { return f.mode }
func (f *fileInfo) ModTime() time.Time { return f.modTime }
func (f *fileInfo) IsDir() bool        { return f.mode.IsDir() }
func (f *fileInfo) Sys() interface{}   { return nil }

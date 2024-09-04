package s3fs

import (
	"io/fs"
	"path"

	"github.com/go-git/go-billy/v5"
)

type S3FS struct{}

func NewS3FS() (billy.Filesystem, error) {
	return &S3FS{}, nil
}

// Chroot implements billy.Filesystem.
func (s *S3FS) Chroot(path string) (billy.Filesystem, error) {
	panic("unimplemented")
}

// Create implements billy.Filesystem.
func (s *S3FS) Create(filename string) (billy.File, error) {
	panic("unimplemented")
}

// Join combines any number of path elements into a single path,
// adding a separator if necessary.
func (s *S3FS) Join(elem ...string) string {
	return path.Join(elem...)
}

// Lstat implements billy.Filesystem.
func (s *S3FS) Lstat(filename string) (fs.FileInfo, error) {
	panic("unimplemented")
}

// MkdirAll implements billy.Filesystem.
func (s *S3FS) MkdirAll(filename string, perm fs.FileMode) error {
	panic("unimplemented")
}

// Open implements billy.Filesystem.
func (s *S3FS) Open(filename string) (billy.File, error) {
	panic("unimplemented")
}

// OpenFile implements billy.Filesystem.
func (s *S3FS) OpenFile(filename string, flag int, perm fs.FileMode) (billy.File, error) {
	panic("unimplemented")
}

// ReadDir implements billy.Filesystem.
func (s *S3FS) ReadDir(path string) ([]fs.FileInfo, error) {
	panic("unimplemented")
}

// Readlink implements billy.Filesystem.
func (s *S3FS) Readlink(link string) (string, error) {
	panic("unimplemented")
}

// Remove implements billy.Filesystem.
func (s *S3FS) Remove(filename string) error {
	panic("unimplemented")
}

// Rename implements billy.Filesystem.
func (s *S3FS) Rename(oldpath string, newpath string) error {
	panic("unimplemented")
}

// Root implements billy.Filesystem.
func (s *S3FS) Root() string {
	panic("unimplemented")
}

// Stat implements billy.Filesystem.
func (s *S3FS) Stat(filename string) (fs.FileInfo, error) {
	panic("unimplemented")
}

// Symlink implements billy.Filesystem.
func (s *S3FS) Symlink(target string, link string) error {
	panic("unimplemented")
}

// TempFile implements billy.Filesystem.
func (s *S3FS) TempFile(dir string, prefix string) (billy.File, error) {
	panic("unimplemented")
}

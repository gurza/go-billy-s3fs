package s3fs

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/go-git/go-billy/v5"
)

type file struct {
	reader  *bytes.Reader
	content []byte
	name    string
}

func newFile(name string, b []byte) (billy.File, error) {
	// use an in-memory reader to wrap the content
	rdr := bytes.NewReader(b)

	return &file{
		reader:  rdr,
		content: b,
		name:    name,
	}, nil
}

func (f *file) Read(b []byte) (int, error) {
	return f.reader.Read(b)
}

func (f *file) ReadAt(b []byte, off int64) (int, error) {
	if off < 0 || off >= int64(len(f.content)) {
		return 0, io.EOF
	}

	slc := f.content[off:]
	n := copy(b, slc)
	if n < len(slc) {
		return n, nil
	}

	return n, io.EOF
}

func (f *file) Write(b []byte) (int, error) {
	f.content = append(f.content, b...)
	f.resetReader()
	return len(b), nil
}

func (f *file) Truncate(size int64) error {
	if size < 0 {
		return &os.PathError{
			Op:   "truncate",
			Path: f.name,
			Err:  errors.New("file size cannot be negative"),
		}
	}

	if size > int64(len(f.content)) {
		padding := make([]byte, size-int64(len(f.content)))
		f.content = append(f.content, padding...)
	} else {
		f.content = f.content[:size]
	}

	f.resetReader()
	return nil
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	return f.reader.Seek(offset, whence)
}

func (f *file) Close() error {
	return nil // no actual file to close
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Size() int64 {
	return int64(len(f.content))
}

func (f *file) Lock() error {
	return ErrLockNotSupported
}

func (f *file) Unlock() error {
	return ErrLockNotSupported
}

func (f *file) resetReader() {
	f.reader = bytes.NewReader(f.content)
}

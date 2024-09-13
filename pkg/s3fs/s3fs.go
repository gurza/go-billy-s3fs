package s3fs

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-git/go-billy/v5"
)

type S3FS struct {
	client *s3.S3
	bucket string
	root   string
}

// NewS3FS creates a new S3-backed filesystem for the given bucket.
func New(client *s3.S3, bucket string) (billy.Filesystem, error) {
	if client == nil {
		return nil, fmt.Errorf("s3 client cannot be nil")
	}
	if bucket == "" {
		return nil, fmt.Errorf("bucket name cannot be empty")
	}

	return &S3FS{
		client: client,
		bucket: bucket,
		root:   "/",
	}, nil
}

// Chroot creates a new filesystem rooted at newRoot within the current root.
func (s *S3FS) Chroot(newRoot string) (billy.Filesystem, error) {
	cleanRoot := path.Clean(newRoot)

	// Ensure the path is relative
	if path.IsAbs(cleanRoot) {
		cleanRoot = cleanRoot[1:]
	}

	newPath := path.Join(s.root, cleanRoot)

	// Ensure the new root path does not escape the current root
	base := path.Clean(s.root) + "/"
	target := path.Clean(newPath) + "/"
	if !strings.HasPrefix(target, base) {
		return nil, fmt.Errorf("invalid path: %s escapes from root", newRoot)
	}

	return &S3FS{
		root: newPath,
	}, nil
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

// Open implements billy.Filesystem.
func (s *S3FS) Open(filename string) (billy.File, error) {
	panic("unimplemented")
}

// OpenFile implements billy.Filesystem.
func (s *S3FS) OpenFile(filename string, flag int, perm fs.FileMode) (billy.File, error) {
	panic("unimplemented")
}

// ReadDir lists the contents of a directory in the S3 bucket,
// returning file and directory information.
func (s *S3FS) ReadDir(name string) ([]os.FileInfo, error) {
	cleanName := path.Clean(name)
	if path.IsAbs(cleanName) {
		cleanName = cleanName[1:]
	}

	// Combine with the root path and normalize
	s3Path := path.Join(s.root, cleanName)
	s3Path = strings.TrimLeft(s3Path, "/")
	if s3Path != "" && !strings.HasSuffix(s3Path, "/") {
		s3Path += "/"
	}

	// Ensure the path doesn't escape the root
	base := path.Clean(strings.TrimLeft(s.root, "/"))
	target := path.Clean(strings.TrimLeft(path.Join(s.root, cleanName), "/"))
	if !isSubPath(base, target) {
		return nil, fmt.Errorf("invalid path: %s escapes from root", name)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.bucket),
		Delimiter: aws.String("/"),
	}
	if s3Path != "" && s3Path != "/" {
		input.Prefix = aws.String(s3Path)
	}

	var results []os.FileInfo
	err := s.client.ListObjectsV2PagesWithContext(ctx, input, func(page *s3.ListObjectsV2Output, last bool) bool {
		for _, prefix := range page.CommonPrefixes {
			dirName := strings.TrimPrefix(*prefix.Prefix, s3Path)
			dirName = strings.TrimSuffix(dirName, "/")
			if dirName != "" && dirName != "/" {
				results = append(results, newDirInfo(dirName))
			}
		}
		for _, obj := range page.Contents {
			fileName := strings.TrimPrefix(*obj.Key, s3Path)
			if fileName != "" && !strings.HasSuffix(fileName, "/") {
				results = append(results, newFileInfo(
					fileName,
					*obj.Size,
					*obj.LastModified,
				))
			}
		}
		return !last
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	return results, nil
}

// MkdirAll creates a directory and all necessary parent directories
// in the S3 bucket with the specified permissions.
func (s *S3FS) MkdirAll(path string, perm fs.FileMode) error {
	return ErrNotImplemented
}

// Lstat implements billy.Filesystem.
func (s *S3FS) Lstat(name string) (fs.FileInfo, error) {
	panic("unimplemented")
}

// Symlink creates newname as a symbolic link to oldname in the S3 bucket.
func (s *S3FS) Symlink(oldname string, newname string) error {
	return ErrNotImplemented
}

// Readlink returns the destination of the named symbolic link
// in the S3 bucket.
func (s *S3FS) Readlink(name string) (string, error) {
	return "", ErrNotImplemented
}

// Remove implements billy.Filesystem.
func (s *S3FS) Remove(filename string) error {
	panic("unimplemented")
}

// Rename implements billy.Filesystem.
func (s *S3FS) Rename(oldpath string, newpath string) error {
	panic("unimplemented")
}

// Root returns the root path of the filesystem.
func (s *S3FS) Root() string {
	return s.root
}

// Stat implements billy.Filesystem.
func (s *S3FS) Stat(filename string) (fs.FileInfo, error) {
	panic("unimplemented")
}

// TempFile implements billy.Filesystem.
func (s *S3FS) TempFile(dir string, prefix string) (billy.File, error) {
	panic("unimplemented")
}

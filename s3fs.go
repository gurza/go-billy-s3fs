package s3fs

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	securejoin "github.com/cyphar/filepath-securejoin"
	"github.com/go-git/go-billy/v5"
)

const SupportedOFlags = os.O_RDONLY

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

// abs converts a filename to an absolute path relative to fs.root,
// ensuring it doesn't escape the root directory.
//
// Symlinks are attempted to be resolved to their target paths. The resolution
// depends on the implementation of the `Lstat` and `Readlink` methods.
// For example, if fs.root is "/base" and fn is "lnk1/lnk2/file.go",
// where "lnk1" is a symlink to "/tgt1" and "lnk2" is a symlink to "/tgt2",
// the result will be "/base/tgt1/tgt2/file".
func (fs *S3FS) abs(fn string) (string, error) {
	return securejoin.SecureJoinVFS(fs.root, fn, fs)
}

// Create implements billy.Filesystem.
func (fs *S3FS) Create(name string) (billy.File, error) {
	return fs.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

// Open implements billy.Filesystem.
func (fs *S3FS) Open(name string) (billy.File, error) {
	return fs.OpenFile(name, os.O_RDONLY, 0)
}

// OpenFile implements billy.Filesystem.
func (fs *S3FS) OpenFile(name string, flag int, perm os.FileMode) (billy.File, error) {
	if flag&SupportedOFlags != flag {
		// todo: support all flags
		return nil, fmt.Errorf("%w: unsupported OpenFile flag %d", ErrNotImplemented, flag)
	}

	resName, err := fs.underlyingPath(name)
	if err != nil {
		return nil, err
	}

	b, err := fs.readObject(resName)
	if err != nil {
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  err,
		}

	}

	f, err := newFile(name, b)
	if err != nil {
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  err,
		}
	}
	return f, nil
}

// readObject retrieves the object content from S3.
func (fs *S3FS) readObject(key string) ([]byte, error) {
	resp, err := fs.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchKey") {
			return nil, os.ErrNotExist
		}
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Join combines any number of path elements into a single path,
// adding a separator if necessary.
func (fs *S3FS) Join(elem ...string) string {
	return path.Join(elem...)
}

// Remove implements billy.Filesystem.
func (fs *S3FS) Remove(filename string) error {
	panic("unimplemented")
}

// Rename implements billy.Filesystem.
func (fs *S3FS) Rename(oldpath string, newpath string) error {
	panic("unimplemented")
}

// Stat retrieves the FileInfo for the named file or directory.
func (fs *S3FS) Stat(name string) (os.FileInfo, error) {
	resName, err := fs.underlyingPath(name)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	input := &s3.HeadObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(resName),
	}
	output, err := fs.client.HeadObjectWithContext(ctx, input)
	if err != nil {
		return nil, &os.PathError{
			Op:   "stat",
			Path: name,
			Err:  fmt.Errorf("no such file or directory"),
		}
	}
	if _, isSymlink := output.Metadata["Symlink-Target"]; isSymlink {
		return nil, fmt.Errorf("%w: symlink handling in Stat()", ErrNotImplemented)
	}

	if strings.HasSuffix(name, "/") {
		return newDirInfo(path.Base(name)), nil
	}
	return newFileInfo(path.Base(name), *output.ContentLength, *output.LastModified), nil
}

// TempFile creates a new temporary file in the directory dir with a name
// beginning with prefix, opens the file for reading and writing, and
// returns the resulting *os.File. If dir is the empty string, TempFile
// uses the default directory for temporary files (see os.TempDir).
// Multiple programs calling TempFile simultaneously will not choose the
// same file. The caller can use f.Name() to find the pathname of the file.
// It is the caller's responsibility to remove the file when no longer
// needed.
func (fs *S3FS) TempFile(dir, prefix string) (billy.File, error) {
	if dir == "" {
		// TODO: implement fs.TempDir()
		// dir = fs.TempDir()
		return nil, fmt.Errorf("%w: TempDir(), prefix = %q", ErrNotImplemented, prefix)
	}

	return nil, ErrNotImplemented
}

// ReadDir lists the contents of a directory in the S3 bucket,
// returning file and directory information.
func (fs *S3FS) ReadDir(name string) ([]os.FileInfo, error) {
	cleanName := path.Clean(name)
	if path.IsAbs(cleanName) {
		cleanName = cleanName[1:]
	}

	// Combine with the root path and normalize
	s3Path := path.Join(fs.root, cleanName)
	s3Path = strings.TrimLeft(s3Path, "/")
	if s3Path != "" && !strings.HasSuffix(s3Path, "/") {
		s3Path += "/"
	}

	// Ensure the path doesn't escape the root
	base := path.Clean(strings.TrimLeft(fs.root, "/"))
	target := path.Clean(strings.TrimLeft(path.Join(fs.root, cleanName), "/"))
	if !isSubPath(base, target) {
		return nil, fmt.Errorf("invalid path: %s escapes from root", name)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(fs.bucket),
		Delimiter: aws.String("/"),
	}
	if s3Path != "" && s3Path != "/" {
		input.Prefix = aws.String(s3Path)
	}

	var results []os.FileInfo
	err := fs.client.ListObjectsV2PagesWithContext(ctx, input, func(page *s3.ListObjectsV2Output, last bool) bool {
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
// within the S3 bucket. Permissions (perm) are ignored.
func (fs *S3FS) MkdirAll(name string, perm os.FileMode) error {
	resName, err := fs.underlyingPath(name)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(resName, "/") {
		// Ensure the path ends with a trailing slash to indicate a "directory"
		resName += "/"
	}
	_, err = fs.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(resName),
		Body:   aws.ReadSeekCloser(strings.NewReader("")),
	})
	if err != nil {
		return &os.PathError{
			Op:   "mkdir",
			Path: name,
			Err:  fmt.Errorf("failed to create directory in S3 bucket %q: %w", fs.bucket, err),
		}
	}

	return nil
}

// Lstat retrieves the FileInfo for the named file or directory
// without following symbolic links.
func (fs *S3FS) Lstat(name string) (os.FileInfo, error) {
	resName, err := fs.underlyingPath(name)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	input := &s3.HeadObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(resName),
	}
	output, err := fs.client.HeadObjectWithContext(ctx, input)
	if err != nil {
		return nil, &os.PathError{
			Op:   "lstat",
			Path: name,
			Err:  fmt.Errorf("no such file or directory"),
		}
	}

	if strings.HasSuffix(name, "/") {
		return newDirInfo(path.Base(name)), nil
	}
	return newFileInfo(path.Base(name), *output.ContentLength, *output.LastModified), nil
}

// Symlink creates newname as a symbolic link to oldname in the S3 bucket.
func (fs *S3FS) Symlink(oldname string, newname string) error {
	return fmt.Errorf("%w: Symlink()", ErrNotImplemented)
}

// Readlink returns the destination of the named symbolic link
// in the S3 bucket.
func (fs *S3FS) Readlink(name string) (string, error) {
	return "", fmt.Errorf("%w: Readlink()", ErrNotImplemented)
}

// Chroot scopes the S3FS to a subdirectory and returns a new S3FS instance
// rooted at the given path.
func (fs *S3FS) Chroot(subPath string) (billy.Filesystem, error) {
	newRoot, err := fs.underlyingPath(subPath)
	if err != nil {
		return nil, err
	}

	return &S3FS{
		client: fs.client,
		bucket: fs.bucket,
		root:   newRoot,
	}, nil
}

// Root returns the root path of the filesystem.
func (fs *S3FS) Root() string {
	return fs.root
}

// underlyingPath ensures the given path is within the allowed boundaries
// and resolves it relative to the current root.
func (fs *S3FS) underlyingPath(p string) (string, error) {
	if isCrossBoundaries(p) {
		return "", billy.ErrCrossedBoundary
	}
	return path.Join(fs.root, path.Clean(p)), nil
}

// isCrossBoundaries checks if the given S3 path escapes boundaries.
func isCrossBoundaries(p string) bool {
	p1 := path.Clean(p)
	return strings.HasPrefix(p1, "../")
}

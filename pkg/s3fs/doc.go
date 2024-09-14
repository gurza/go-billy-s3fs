// Package s3fs provides an implementation of the billy.Filesystem interface
// for interacting with an object storage bucket (e.g. Amazon S3 bucket)
// as if it were a traditional file system. This allows seamless integration
// with applications that rely on the billy abstraction, enabling operations
// such as reading, writing, deleting, and traversing files and directories.
package s3fs

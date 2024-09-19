// Package s3fs implements the billy.Filesystem interface for interacting
// with an Amazon S3 bucket as a virtual filesystem. This enables seamless
// integration with applications relying on the billy abstraction, supporting
// operations like reading, writing, deleting, and traversing files and dirs.
//
// Note: When working with paths in this package, always use the `path` package
// instead of `filepath`, as S3 paths are UNIX-like and must use forward
// slashes (/).
package s3fs

package s3fs

import "errors"

var (
	ErrLockNotSupported = errors.New("locking is not supported")
	ErrNotImplemented   = errors.New("not implemented")
)

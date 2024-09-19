package s3fs

import (
	"reflect"
	"testing"

	"github.com/go-git/go-billy/v5"
)

func TestS3FS_ImplementsBillyFilesystem(t *testing.T) {
	var fsys S3FS

	iface := reflect.TypeOf((*billy.Filesystem)(nil)).Elem()
	sType := reflect.TypeOf(&fsys)

	// Check if *S3FS implements billy.Filesystem
	if !sType.Implements(iface) {
		t.Errorf("S3FS does not implement billy.Filesystem interface")
	}
}

package s3fs

import (
	"reflect"
	"testing"

	"github.com/go-git/go-billy/v5"
)

func TestS3FS_ImplementsBillyFile(t *testing.T) {
	var f file

	iface := reflect.TypeOf((*billy.File)(nil)).Elem()
	sType := reflect.TypeOf(&f)

	// Check if *file implements billy.File
	if !sType.Implements(iface) {
		t.Errorf("file does not implement billy.File interface")
	}
}

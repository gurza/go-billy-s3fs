package s3fs

import (
	"io/fs"
	"reflect"
	"testing"
)

func TestFileInfo_ImplementsFSFileInfo(t *testing.T) {
	var fi fileInfo

	iface := reflect.TypeOf((*fs.FileInfo)(nil)).Elem()
	sType := reflect.TypeOf(&fi)

	// Check if *fileInfo implements fs.FileInfo
	if !sType.Implements(iface) {
		t.Errorf("fileInfo does not implement fs.FileInfo interface")
	}
}

package s3fs

import (
	"os"
	"reflect"
	"testing"
)

func TestFileStat_ImplementsOSFileInfo(t *testing.T) {
	var fs fileStat

	iface := reflect.TypeOf((*os.FileInfo)(nil)).Elem()
	sType := reflect.TypeOf(&fs)

	if !sType.Implements(iface) {
		t.Errorf("fileStat does not implement os.FileInfo interface")
	}
}

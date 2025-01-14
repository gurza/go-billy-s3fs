package s3fs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSubPath(t *testing.T) {
	tests := []struct {
		basepath string
		targpath string
		want     bool
	}{
		// Valid subpaths within the basepath
		{"foo", "foo/bar", true},
		{"foo/", "foo/bar/", true},
		{"/", "foo/bar", true}, // basepath is root, allowing any subpath
		{".", "foo/bar", true}, // basepath is current directory
		{"", "foo", true},      // basepath is empty, treated as root

		// Preventing directory traversal outside the basepath
		{"foo", "foobar", false},
		{"foo", "foo/../../baz", false},            // traverses outside basepath
		{"foo/bar", "foo/bar/../../../baz", false}, // multiple traversals outside basepath
		{"foo", "../foo/bar", false},               // absolute traversal outside basepath

		// Exact matches
		{"foo", "foo", true},
		{"foo/bar", "foo/bar", true},
		{"/", "/", true},
		{".", ".", true},
		{"", "", true},

		// Edge cases related to normalization
		{"foo", "foo/bar/..", true},
		{"foo", "foo/bar/baz/..", true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("base=%s,target=%s", tt.basepath, tt.targpath), func(t *testing.T) {
			got := isSubPath(tt.basepath, tt.targpath)
			if got != tt.want {
				t.Errorf("isSubPath(%q, %q) = %v; want %v", tt.basepath, tt.targpath, got, tt.want)
			}
		})
	}
}

func TestLastIndexByte(t *testing.T) {
	// byte is found at the beginning
	assert.Equal(t, 0, lastIndexByte("abcdef", 'a'))

	// byte is found in the middle
	assert.Equal(t, 2, lastIndexByte("abcdef", 'c'))

	// byte is found at the end
	assert.Equal(t, 5, lastIndexByte("abcdef", 'f'))

	// byte is found multiple times (should return last occurrence)
	assert.Equal(t, 6, lastIndexByte("abcdefcd", 'c'))

	// byte is not found
	assert.Equal(t, -1, lastIndexByte("abcdef", 'z'))

	// empty string
	assert.Equal(t, -1, lastIndexByte("", 'a'))
}

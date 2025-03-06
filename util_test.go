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

func TestPrefixAndSuffix(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		prefix  string
		suffix  string
	}{
		{
			name:    "no wildcard",
			pattern: "abc",
			prefix:  "abc",
			suffix:  "",
		},
		{
			name:    "wildcard at the beginning",
			pattern: "*def",
			prefix:  "",
			suffix:  "def",
		},
		{
			name:    "wildcard in the middle",
			pattern: "abc*def",
			prefix:  "abc",
			suffix:  "def",
		},
		{
			name:    "wildcard at the end",
			pattern: "abc*",
			prefix:  "abc",
			suffix:  "",
		},
		{
			name:    "multiple wildcards (split by last one)",
			pattern: "abc*def*ghi",
			prefix:  "abc*def",
			suffix:  "ghi",
		},
		{
			name:    "empty string",
			pattern: "",
			prefix:  "",
			suffix:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefix, suffix, err := prefixAndSuffix(tt.pattern)
			assert.NoError(t, err)
			assert.Equal(t, tt.prefix, prefix)
			assert.Equal(t, tt.suffix, suffix)
		})
	}
}

func TestPrefixAndSuffix_Error(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
	}{
		{
			name:    "path separator",
			pattern: "abc/def",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefix, suffix, err := prefixAndSuffix(tt.pattern)
			assert.Error(t, err)
			assert.Equal(t, "", prefix)
			assert.Equal(t, "", suffix)
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

func TestJoinPath(t *testing.T) {
	tests := []struct {
		name string
		dir  string
		fn   string
		want string
	}{
		{
			name: "directory with separator",
			dir:  "/home/user/",
			fn:   "file.txt",
			want: "/home/user/file.txt",
		},
		{
			name: "directory without separator",
			dir:  "/home/user",
			fn:   "file.txt",
			want: "/home/user/file.txt",
		},
		{
			name: "empty directory",
			dir:  "",
			fn:   "file.txt",
			want: "/file.txt",
		},
		{
			name: "directory with separator and empty filename",
			dir:  "/home/user/",
			fn:   "",
			want: "/home/user/",
		},
		{
			name: "directory without separator and empty filename",
			dir:  "/home/user",
			fn:   "",
			want: "/home/user/",
		},
		{
			name: "empty directory and filename",
			dir:  "",
			fn:   "",
			want: "/",
		},
		{
			name: "relative path",
			dir:  "user/docs",
			fn:   "file.txt",
			want: "user/docs/file.txt",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := joinPath(tc.dir, tc.fn)
			assert.Equal(t, tc.want, got)
		})
	}
}

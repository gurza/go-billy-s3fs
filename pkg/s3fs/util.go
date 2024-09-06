package s3fs

import (
	"path"
	"path/filepath"
	"strings"
)

// isSubPath checks if targpath is a subpath of basepath.
// It ensures that targpath resides within basepath, thereby preventing
// directory traversal outside the basepath.
func isSubPath(basepath, targpath string) bool {
	basepath = path.Clean(basepath)
	targpath = path.Clean(targpath)

	// Remove leading slashes to treat paths as relative
	basepath = strings.TrimPrefix(basepath, "/")
	targpath = strings.TrimPrefix(targpath, "/")

	if basepath == "." || basepath == "" {
		// If basepath is root, any targpath is allowed
		return true
	}

	rel, err := filepath.Rel(basepath, targpath)
	if err != nil {
		return false
	}
	if strings.HasPrefix(rel, "..") {
		// If the relative path starts with "..", targpath is outside basepath
		return false
	}

	// Optionally, prevent targpath from being exactly the basepath
	// Uncomment the following lines if you want to disallow exact matches
	/*
	   if rel == "." {
	       return false
	   }
	*/

	return true
}

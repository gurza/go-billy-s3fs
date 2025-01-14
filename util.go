package s3fs

import (
	"errors"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
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

// prefixAndSuffix splits pattern by the last wildcard "*", if applicable,
// returning prefix as the part before "*" and suffix as the part after "*".
func prefixAndSuffix(pattern string) (prefix, suffix string, err error) {
	// FIXME: use range after with golang122
	for i := 0; i < len(pattern); i++ {
		if pattern[i] == PathSeparator {
			return "", "", errors.New("pattern contains path separator")
		}
	}
	if pos := lastIndexByte(pattern, '*'); pos != -1 {
		prefix, suffix = pattern[:pos], pattern[pos+1:]
	} else {
		prefix = pattern
	}
	return prefix, suffix, nil
}

func lastIndexByte(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func getRandom() string {
	uuid := uuid.New().String()
	return strings.ReplaceAll(uuid, "-", "")
}

func joinPath(dir, fn string) string {
	if len(dir) > 0 && dir[len(dir)-1] == PathSeparator {
		return dir + fn
	}
	return dir + string(PathSeparator) + fn
}

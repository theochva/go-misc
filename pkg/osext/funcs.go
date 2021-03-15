package osext

import (
	"io/fs"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// ReadFile - read the contents of a file and optionally trim any whitespace.
func ReadFile(filename string, trim bool) (bytes []byte, err error) {
	if bytes, err = os.ReadFile(filename); err != nil {
		return nil, errors.Wrapf(err, "Failed to read file '%s'", filename)
	}

	if trim {
		contents := strings.TrimSpace(string(bytes))
		bytes = []byte(contents)
	}

	return
}

// ReadFileAsString - read the contents of a file and optionally trim any whitespace.
func ReadFileAsString(filename string, trim bool) (contents string, err error) {
	var bytes []byte

	if bytes, err = os.ReadFile(filename); err != nil {
		return "", errors.Wrapf(err, "Failed to read file '%s'", filename)
	}

	contents = string(bytes)

	if trim {
		contents = strings.TrimSpace(contents)
	}

	return
}

// CreateTempWithContents - create a temp file with the filename pattern and contents specified.
//
// It uses os.CreateTemp(dir, filenamePattern) to create the file and then just writes the contents provided.
func CreateTempWithContents(dir, fileNamePattern string, bytes []byte, perm fs.FileMode) (file *os.File, err error) {
	if file, err = os.CreateTemp(dir, fileNamePattern); err != nil {
		return nil, err
	}

	if len(bytes) > 0 {
		if err = os.WriteFile(file.Name(), bytes, perm); err != nil {
			return nil, errors.Wrapf(err, "Failed to write contents to temp file '%s'", file.Name())
		}
	}

	return
}

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	if filename == "" {
		return false
	}

	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// FileOrDirectoryExists - check if specified filename (which is a file or directory) exists
func FileOrDirectoryExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// IsDirectory - check if specified filename is a directory
func IsDirectory(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

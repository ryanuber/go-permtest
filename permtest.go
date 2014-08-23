package permtest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Write determines if it is possible to write to the file or directory
// indicated by <path>. Special files like pipes, sockets, and devices will
// return an error.
func Write(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("%s: permission denied", path)
		}

		if os.IsNotExist(err) {
			return Write(filepath.Dir(path))
		}

		return err
	}

	switch {
	case fi.Mode().IsRegular():
		return WriteFile(path)
	case fi.Mode().IsDir():
		return WriteDir(path)
	default:
		return fmt.Errorf("%s: not a file or directory", path)
	}
}

// WriteFile determines if a file is able to be written to by attempting to open
// it for writing. If the file does not exist, then its parent directory will be
// checked for write permission instead. This will continue to traverse up the
// requested directory structure until a writable directory or error is
// encountered.
func WriteFile(path string) error {
	fh, err := os.OpenFile(path, os.O_APPEND, 0666)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("%s: permission denied", path)
		}

		if os.IsNotExist(err) {
			return WriteDir(filepath.Dir(path))
		}
	}
	defer fh.Close()
	return nil
}

// WriteDir checks if a given directory is writable or not. If the directory
// path does not exist, its parent will be checked. This will continue until
// either a writable directory is found, or an error is encountered.
//
// Directories are tested by writing a temporary hidden file into them. This
// file will be removed immediately after the test.
func WriteDir(path string) error {
	fh, err := ioutil.TempFile(path, ".permtest-")
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("%s: permission denied", path)
		}
		if os.IsNotExist(err) {
			return WriteDir(filepath.Dir(path))
		}
		return err
	}
	defer fh.Close()
	defer os.Remove(fh.Name())

	return nil
}

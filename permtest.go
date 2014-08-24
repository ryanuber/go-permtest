package permtest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const errDenied string = "%s: permission denied"

// Write determines if it is possible to write to the file or directory
// indicated by <path>. Special files like pipes, sockets, and devices will
// return an error.
//
// This function returns two values: The last path which was write-tested, and
// an error (or nil). Since the function recurses, having the last tested path
// returned can be helpful in determining root cause on where access was
// actually denied.
func Write(path string) (string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		switch {
		case os.IsPermission(err):
			return path, fmt.Errorf(errDenied, path)
		case os.IsNotExist(err):
			return Write(filepath.Dir(path))
		default:
			return path, err
		}
	}

	switch {
	case fi.Mode().IsRegular():
		return WriteFile(path)
	case fi.Mode().IsDir():
		return WriteDir(path)
	default:
		return path, fmt.Errorf("%s: not a file or directory", path)
	}
}

// WriteFile determines if a file is able to be written to by attempting to open
// it for writing. If the file does not exist, an error will be thrown.
//
// This function errors on non-existent files because it would be too ambiguous
// to only know that it was possible, if parent directories were first created,
// to create a new file called <path>. Instead, we will assume here that the
// file already exists, and that the user just wants to know if they are able
// to write to it.
func WriteFile(path string) (string, error) {
	fh, err := os.OpenFile(path, os.O_APPEND, 0666)
	if err != nil {
		switch {
		case os.IsPermission(err):
			return path, fmt.Errorf(errDenied, path)
		case os.IsNotExist(err):
			if err := writeTempFile(filepath.Dir(path)); err != nil {
				if os.IsPermission(err) {
					return path, fmt.Errorf(errDenied, path)
				}
				return path, err
			}
			return path, nil
		default:
			return path, err
		}
	}
	defer fh.Close()
	return path, nil
}

// WriteDir checks if a given directory is writable or not. If the directory
// path does not exist, its parent will be checked. This will continue until
// either a writable directory is found, or an error is encountered.
func WriteDir(path string) (string, error) {
	if err := writeTempFile(path); err != nil {
		switch {
		case os.IsPermission(err):
			return path, fmt.Errorf(errDenied, path)
		case os.IsNotExist(err):
			return WriteDir(filepath.Dir(path))
		default:
			return path, err
		}
	}
	return path, nil
}

// writeTempFile will attempt to write a temporary file into a directory. This
// determines beyond any shred of doubt whether or not the user is able to write
// a file to a given path. The file will be automatically deleted as soon as
// this method returns, and the file should have a unique, random name so as to
// not interfere with any existing content.
func writeTempFile(dir string) error {
	fh, err := ioutil.TempFile(dir, ".permtest-")
	if err != nil {
		return err
	}
	defer fh.Close()
	defer os.Remove(fh.Name())
	return nil
}

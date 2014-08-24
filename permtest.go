package permtest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const errDenied string = "%s: permission denied"

// Write determines if it is possible to write to the file or directory
// indicated by <path>.
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

	if fi.IsDir() {
		return WriteDir(path)
	}

	return WriteFile(path)
}

// WriteFile determines if a file is able to be written to by attempting to open
// it for writing. If the file does not exist, then its parent directory
// is tested for write capability by the current user.
//
// This method will return an error if one of two conditions are met:
// 1. The file exists and is not writable
// 2. The file does not exist and its parent directory is not writable
//
// The intent is that after testing with this method, a call to os.Create()
// would be made to open or create the file and write some data to it.
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

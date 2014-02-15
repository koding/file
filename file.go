// Package file provides simple utility functions to be used with files
package file

import (
	"errors"
	"io"
	"os"
)

// FileInfo describes a file. It's a wrapper around os.FileInfo.
type FileInfo struct {
	Exists bool
	os.FileInfo
}

// Stat returns a FileInfo describing the named file. It's a wrapper around os.Stat.
func Stat(file string) (*FileInfo, error) {
	fi, err := os.Stat(file)
	if err == nil {
		f := &FileInfo{Exists: true}
		f.FileInfo = fi
		return f, nil
	}

	if os.IsNotExist(err) {
		f := &FileInfo{Exists: false}
		f.FileInfo = nil
		return f, nil
	}

	return nil, err
}

// IsFile checks wether the given file is a directory or not. It panics if an
// error is occured. Use IsFileOk to use the returned error.
func IsFile(file string) bool {
	ok, err := IsFileOk(file)
	if err != nil {
		panic(err)
	}

	return ok
}

// IsFileOk checks whether the given file is a directory or not.
func IsFileOk(file string) (bool, error) {
	sf, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer sf.Close()

	fi, err := sf.Stat()
	if err != nil {
		return false, err
	}

	if fi.IsDir() {
		return false, nil
	}

	return true, nil
}

// Exists checks whether the given file exists or not. It panics if an error
// is occured. Use ExistsOk to use the returned error.
func Exists(file string) bool {
	ok, err := ExistsOk(file)
	if err != nil {
		panic(err)
	}

	return ok
}

// ExistsOk checks whether the given file exists or not.
func ExistsOk(file string) (bool, error) {
	_, err := os.Stat(file)
	if err == nil {
		return true, nil // file exist
	}

	if os.IsNotExist(err) {
		return false, nil // file does not exist
	}

	return false, err
}

func copyFile(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	fi, err := sf.Stat()
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return errors.New("src is a directory, please provide a file")
	}

	df, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fi.Mode())
	if err != nil {
		return err
	}
	defer df.Close()

	if _, err := io.Copy(df, sf); err != nil {
		return err
	}

	return nil
}

// TODO: implement those functions
func isReadable(mode os.FileMode) bool { return mode&0400 != 0 }

func isWritable(mode os.FileMode) bool { return mode&0200 != 0 }

func isExecutable(mode os.FileMode) bool { return mode&0100 != 0 }

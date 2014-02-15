// Package file provides simple utility functions to be used with files
package file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Copy copies the file or directory from source path to destination path.
// Directories are copied recursively. Copy does not handle symlinks currently.
func Copy(src, dst string) error {
	if dst == "." {
		dst = filepath.Base(src)
	}

	if src == dst {
		return fmt.Errorf("%s and %s are identical (not copied).", src, dst)
	}

	if !Exists(src) {
		return fmt.Errorf("%s: no such file or directory.", src)
	}

	if Exists(dst) && IsFile(dst) {
		return fmt.Errorf("%s is a directory (not copied).", src)
	}

	srcBase, _ := filepath.Split(src)
	walks := 0

	// dstPath returns the rewritten destination path for the given source path
	dstPath := func(srcPath string) string {
		// some/random/long/path/example/hello.txt -> example/hello.txt
		srcPath = strings.TrimPrefix(srcPath, srcBase)

		// example/hello.txt -> destination/example/hello.txt
		if walks != 0 {
			return filepath.Join(dst, srcPath)
		}

		// hello.txt -> example/hello.txt
		if Exists(dst) && !IsFile(dst) {
			return filepath.Join(dst, filepath.Base(srcPath))
		}

		// hello.txt -> test.txt
		return dst
	}

	filepath.Walk(src, func(srcPath string, file os.FileInfo, err error) error {
		defer func() { walks++ }()

		if file.IsDir() {
			fmt.Printf("copy dir from '%s' to '%s'\n", srcPath, dstPath(srcPath))
			os.MkdirAll(dstPath(srcPath), 0755)
		} else {
			fmt.Printf("copy file from '%s' to '%s'\n", srcPath, dstPath(srcPath))
			err = copyFile(srcPath, dstPath(srcPath))
			if err != nil {
				fmt.Println(err)
			}
		}

		return nil
	})

	return nil
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

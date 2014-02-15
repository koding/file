# File [![GoDoc](https://godoc.org/github.com/koding/file?status.png)](http://godoc.org/github.com/koding/file)

File includes several useful file related helper functions. Especially common
tasks like copying a file or directory, checking if a file exists, etc. are
not easy to write in Go.

For usage see examples below or click on the godoc badge.

## Install

```bash
go get github.com/koding/file
```

## Examples

```go
// copy recursively exampleDir to a new test directory
err := file.Copy("exampleDir", "testDir")

// copy a file into a folder
err := file.Copy("hello.txt", "./exampleDir")

// create a copy of a given file
err := file.Copy("hello.txt", "another.txt")

// check if a file exists
ok := file.Exists("hello.txt")

// Is the given path a file?
ok := file.IsFile("another.png")

```

## License

The MIT License (MIT) - see LICENSE.md for more details

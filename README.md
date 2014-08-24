# go-permtest

permtest is a Go library providing os-generic tests for basic file permissions.
This can be useful for testing a user's ability to create files and directories
on a filesystem without leaving a mess behind.

An example use case is testing one's ability to write a file to
`/foo/bar/baz/zip`. If you simply catch errors from `os.MkdirAll()`, you will
not have any idea which directories were or were not created successfully,
meaning you can't clean up after yourself by removing the created directories
and files. `permtest` fixes this by testing write access to each directory from
deepest upward by writing a temporary file and then immediately removing it.

`permtest` is able to determine if a directory structure is able to be created
via its `permtest.WriteDir()` method, or if a file is able to be created using
the `permtest.WriteFile()` method.

Example:

```
package main

import (
    "fmt"
    "github.com/ryanuber/go-permtest"
)

func main() {
    // test if we have the ability to create a file in /tmp
    if _, err := permtest.WriteFile("/tmp/foo"); err != nil {
        fmt.Println(err.Error())
    }

    // test if we can write a file directly into a non-existent directory
    if _, err := permtest.WriteFile("/tmp/foo/bar"); err != nil {
        fmt.Println(err.Error())
    }

    // test is we have the ability to create a deep directory structure
    if _, err := permtest.WriteDir("/tmp/foo/bar/baz/zip"); err != nil {
        fmt.Println(err.Error())
    }
}
```

`permtest` returns both an error as well as the last tested path as a string,
which makes it easy to figure out where permission was actually denied within
the file structure.

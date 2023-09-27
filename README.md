# go-editline

[![PkgGoDev](https://pkg.go.dev/badge/github.com/maolonglong/go-editline)](https://pkg.go.dev/github.com/maolonglong/go-editline)

Go bindings for the [editline](https://github.com/troglobit/editline).

## Installation

This library depends on editline and requires it to be installed beforehand. You can refer to the following link for instructions on how to build and install editline: <https://github.com/troglobit/editline#build--install>.

```bash
go get github.com/maolonglong/go-editline
```

## Usage

Some useful hints on how to use the library is available in the [examples/](./examples/) directory.

```go
package main

import (
	"fmt"
	"io"

	"github.com/maolonglong/go-editline"
)

func main() {
	defer editline.Uninitialize()

	for {
		line, err := editline.ReadLine("> ")
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		fmt.Println(line)
	}
}
```

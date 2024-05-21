# Go Contemplate

[![Build Status](https://github.com/foomo/gocontemplate/actions/workflows/test.yml/badge.svg?branch=main&event=push)](https://github.com/foomo/gocontemplate/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/foomo/gocontemplate)](https://goreportcard.com/report/github.com/foomo/gocontemplate)
[![GoDoc](https://godoc.org/github.com/foomo/gocontemplate?status.svg)](https://godoc.org/github.com/foomo/gocontemplate)

> A code generation helper.

Wrapper library around `golang.org/x/tools/go/packages` to filter only defined types and their dependencies.

## Example

```go
package main

import (
  "github.com/foomo/gocontemplate"
)

func main() {
  goctpl, err := gocontemplate.Load(&gocontemplate.Config{
    Packages: []*gocontemplate.ConfigPackage{
      {
        Path:  "github.com/foomo/sesamy-go/event",
        Types: []string{"PageView"},
      },
    },
  })
  if err != nil {
    panic(err)
  }
}
```

## How to Contribute

Make a pull request...

## License

Distributed under MIT License, please see license file within the code for more details.

[![CodeQL](https://github.com/foomo/gocontemplate/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/foomo/gocontemplate/actions/workflows/github-code-scanning/codeql)
[![Coverage](https://img.shields.io/codecov/c/github/foomo/gocontemplate?style=flat-square&logo=github)](https://app.codecov.io/gh/foomo/gocontemplate)
[![GoDoc](https://godoc.org/github.com/foomo/gocontemplate?status.svg)](https://godoc.org/github.com/foomo/gocontemplate)
[![GitHub Stars](https://img.shields.io/github/stars/foomo/gocontemplate.svg?style=flat-square&logo=github)](https://github.com/foomo/gocontemplate)

<p align="center">
  <img alt="gocontemplate" src="docs/public/logo.png" width="400" height="400"/>
</p>

# Go Contemplate

> A code generation helper.

Wrapper library around `golang.org/x/tools/go/packages` to filter only defined types and their dependencies.

## Example

```go
package main

import (
  "github.com/foomo/gocontemplate/pkg/contemplate"
)

func main() {
  ctpl, err := contemplate.Load(&contemplate.Config{
    Packages: []*contemplate.PackageConfig{
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

Contributions are welcome! Please read the [contributing guide](docs/CONTRIBUTING.md).

![Contributors](https://contributors-table.vercel.app/image?repo=foomo/gocontemplate&width=50&columns=15)

## License

Distributed under MIT License, please read the [license file](LICENSE) for more details.

_Made with ♥ [foomo](https://www.foomo.org) by [bestbytes](https://www.bestbytes.com)_

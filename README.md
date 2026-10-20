[![CodeQL](https://github.com/foomo/gocontemplate/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/foomo/gocontemplate/actions/workflows/github-code-scanning/codeql)
[![Coverage](https://img.shields.io/codecov/c/github/foomo/gocontemplate?style=flat-square&logo=github)](https://app.codecov.io/gh/foomo/gocontemplate)
[![GoDoc](https://pkg.go.dev/badge/github.com/foomo/gocontemplate.svg)](https://pkg.go.dev/github.com/foomo/gocontemplate)
[![GitHub Stars](https://img.shields.io/github/stars/foomo/gocontemplate.svg?style=flat-square&logo=github)](https://github.com/foomo/gocontemplate)

<p align="center">
  <img alt="gocontemplate" src="docs/public/logo.png" width="400" height="400"/>
</p>

# Go Contemplate

> A code generation helper.

Wrapper library around `golang.org/x/tools/go/packages` that loads Go packages and narrows them down to the types you ask for, plus every type those types depend on — across package boundaries.

Typical use: a code generator that needs the full type graph behind a few root types (e.g. to emit TypeScript, JSON Schema or docs) without walking the AST and imports itself.

## Installation

```shell
go get github.com/foomo/gocontemplate
```

## Usage

```go
package main

import (
  "context"
  "fmt"

  "github.com/foomo/gocontemplate/pkg/contemplate"
)

func main() {
  ctpl, err := contemplate.Load(context.Background(), &contemplate.Config{
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

  // Print every type that PageView depends on, grouped by package.
  for path, pkg := range ctpl.Packages {
    for name := range pkg.ScopeTypes() {
      fmt.Println(path, name)
    }
  }
}
```

Given these types:

```go
package test

type Event[T any] struct {
  Name   string `json:"name"`
  Params T      `json:"params,omitempty"`
}
```

```go
package event

type PageView test.Event[params.PageView]
```

requesting `event.PageView` puts `event.PageView`, `test.Event` and `params.PageView` (and everything they reference) in scope.

## Concepts

- **`Contemplate.Packages`** holds every loaded package, including transitive imports, keyed by import path.
- **`Package.Types()`** returns all objects in a package scope.
- **`Package.ScopeTypes()`** returns only the requested types and the types they reference. This is usually what a generator iterates over.
- **`Contemplate.LookupTypesByType()`** finds constants and types derived from a named type, e.g. the values of an enum.
- **`Package.Raw()`** exposes the underlying `*packages.Package` for anything not covered.

## Configuration

`Config` carries `json` and `yaml` tags, so it can be embedded in a generator's config file:

```yaml
directory: ./ # directory containing go.mod; env vars are expanded; defaults to cwd
packages:
  - path: github.com/foomo/sesamy-go/event
    types:
      - PageView
```

See the [package documentation](https://pkg.go.dev/github.com/foomo/gocontemplate/pkg/contemplate) for the full API.

## How to Contribute

Contributions are welcome! Please read the [contributing guide](docs/CONTRIBUTING.md).

![Contributors](https://contributors-table.vercel.app/image?repo=foomo/gocontemplate&width=50&columns=15)

## License

Distributed under MIT License, please read the [license file](LICENSE) for more details.

_Made with ♥ [foomo](https://www.foomo.org) by [bestbytes](https://www.bestbytes.com)_

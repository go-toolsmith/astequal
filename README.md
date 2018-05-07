[![Go Report Card](https://goreportcard.com/badge/github.com/Quasilyte/astcmp)](https://goreportcard.com/report/github.com/Quasilyte/astcmp)
[![GoDoc](https://godoc.org/github.com/Quasilyte/astcmp?status.svg)](https://godoc.org/github.com/Quasilyte/astcmp)

# astcmp

Compare AST nodes by their printed representations (strings) or directly.

## Use cases

`astcmp` is intended to be used by Go tools that examine `go/ast` nodes for
specific patterns.

This list includes programs like:
- Static code analyzers and linters
- Go code analyzers
- Source-level code optimizers


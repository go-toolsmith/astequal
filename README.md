[![Go Report Card](https://goreportcard.com/badge/github.com/Quasilyte/astcmp)](https://goreportcard.com/report/github.com/Quasilyte/astcmp)
[![GoDoc](https://godoc.org/github.com/Quasilyte/astcmp?status.svg)](https://godoc.org/github.com/Quasilyte/astcmp)

# astcmp

Compare AST nodes by their printed representations (strings) or directly.

## Quick start

### Installation:

```
go get -u github.com/Quasilyte/astcmp
```

### Simple example"

```go
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"

	"github.com/Quasilyte/astcmp"
)

func main() {
	const code = `
		package foo

		func main() {
			x := []int{1, 2, 3}
			x := []int{1, 2, 3}
		}`

	fset := token.NewFileSet()
	pkg, err := parser.ParseFile(fset, "string", code, 0)
	if err != nil {
		log.Fatalf("parse error: %+v", err)
	}

	fn := pkg.Decls[0].(*ast.FuncDecl)
	x := fn.Body.List[0]
	y := fn.Body.List[1]
	fmt.Println(reflect.DeepEqual(x, y)) // => false
	fmt.Println(astcmp.Equal(x, y))      // => true
	fmt.Println(astcmp.EqualStmt(x, y))  // => true

	fmt.Println(astcmp.EqualExprString(`x + y`, `x+y`)) // => true
}
```

## Use cases

`astcmp` is intended to be used by Go tools that examine `go/ast` nodes for
specific patterns.

This list includes programs like:
- Go static code analyzers and linters
- Source-level code optimizers


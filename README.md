# astcmp

Compare AST nodes by patterns, strings or directly.

## Use cases

`astcmp` is intended to be used by Go tools that examine `go/ast` nodes for
specific patterns.

This list includes programs like:
- Static code analyzers and linters
- Go code analyzers
- Editor/IDE extensions for code search by syntax-aware queries
- Source-level code optimizers

## Description

TODO.

The best way to learn about pattern matching is to inspect these test suites:
- [Matching expressions](src/astcmp/testdata/simple_expr.txt)
- [Matching statements](src/astcmp/testdata/simple_stmt.txt)
- [Matching declarations](src/astcmp/testdata/simple_decl.txt)
- [Matching complex statements](src/astcmp/testdata/complex_stmt.txt)
- [Matching complex expressions](src/astcmp/testdata/complex_expr.txt)
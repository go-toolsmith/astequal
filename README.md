# astexp
Like regular expressions, but for Go AST.

## Description

This description is taken from `astexp/doc.go` as a temporary placeholder.

```
 Package astexp implements Go AST search.

 This is pre-release, unstable version of library.
 TODO: remove this warning after v1.0 is reached.

 Patterns have syntax of S-expressions.

 Each Go AST node have S-expresion equivalent defined.
 Example patterns:
	`"123"` matches `"123"`
	`(+ 1 2)` matches `1+2`
	`(call f "1")` matches `f("1")`
	`(block (callstmt println "hi"))` matches `{println("hi")}`
 See astexp_test.go for more complex examples.

 Main operations:
	Match -- tell if AST node matches the pattern
	Find -- retrieve matching AST sub-node(s)

 Limitations:
	- ast.ParenExpr is ignored and can't be expressed in patterns.
	- comments are ignored.

 Warning: undocummented/untested pattern combinations
 may behave differently from version to version.
 Weird behavior that is proved to be unintended, will
 be fixed even if it can break someone who uses this.
 Docummented (and tested) features, on the other hand, have
 strong backwards compatibility from sexp v1.0 and above.
 ```

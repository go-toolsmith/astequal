package astexp

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

type matchTest struct {
	input string
	tmpl  string
}

func TestMatchExprAccept(t *testing.T) {
	tests := []matchTest{
		{`1`, `1`},
		{`(1)`, `1`},
		{`((1))`, `1`},
		{`x`, `x`},

		{`1 + 2`, `(+ 1 2)`},
		{`x + 1`, `(+ x 1)`},
		{`x + y + z`, `(+ (+ x y) z)`},
		{`x + (y + z)`, `(+ x (+ y z))`},
		{`x + y * z`, `(+ x (* y z))`},
		{`(x + y) * z`, `(* (+ x y) z)`},

		{`f(1)`, `(call f 1)`},
		{`append(xs, 1)`, `(call append xs 1)`},
		{`append(xs, a, b, 1)`, `(call append xs $... 1)`},
		{`f(1)(2)`, `(call (call f 1) 2)`},

		{`foo.bar`, `(. foo bar)`},
		{`foo.bar.baz`, `(. (. foo bar) baz)`},
		{`foo.bar()`, `(call (. foo bar))`},
		{`foo.bar.baz(x)`, `(call (. (. foo bar) baz) x)`},
		{`a.x + b.y + c.z`, `(+ (+ (. a x) (. b y)) (. c z))`},

		{`f().x`, `(. (call f) x)`},
		{`f().g().x`, `(. (call (. (call f) g)) x)`},

		{`1 + 2.3`, `(+ 1 $float)`},
		{`1 + 2.3`, `(+ $int 2.3)`},
		{`1 + 2.3`, `(+ $int $float)`},
		{`"a" + "b"`, `(+ "a" $str)`},
		{`"a" + "b"`, `(+ $str "b")`},
		{`"a" + "b"`, `(+ $str $str)`},
		{`ch - '0'`, `(- ch $char)`},
		{`ch - '0'`, `(- $id '0')`},
		{`ch - '0'`, `(- $id $char)`},

		{`[5]int`, `(array_type 5 int)`},
		{`[foo + 1][10]float32`, `(array_type (+ foo 1) (array_type 10 float32))`},
		{`[]int`, `(slice_type int)`},
		{`[][]int`, `(slice_type (slice_type int))`},
		{`[][][]Point`, `(slice_type (slice_type (slice_type Point)))`},
	}

	for _, test := range tests {
		runMatchTest(t, true, test, astExpr(t, test.input))
	}
}

func TestMatchExprReject(t *testing.T) {
	tests := []matchTest{
		{`f()`, `(callstmt f)`},
		{`1 + f(x)`, `(+ 1 (callstmt f x))`},
	}

	for _, test := range tests {
		runMatchTest(t, false, test, astExpr(t, test.input))
	}
}

func TestMatchStmtAccept(t *testing.T) {
	tests := []matchTest{
		{`{1; 2}`, `(block 1 2)`},
		{`{f()}`, `(block (call f))`},
		{`{f()}`, `(block (callstmt f))`},

		{`{foo.bar}`, `(block (. foo bar))`},
	}

	for _, test := range tests {
		runMatchTest(t, true, test, astStmt(t, test.input))
	}
}

func TestMatchStmtReject(t *testing.T) {
	tests := []matchTest{
		{`{1}`, `1`},
	}

	for _, test := range tests {
		runMatchTest(t, false, test, astStmt(t, test.input))
	}
}

func runMatchTest(t *testing.T, shouldAccept bool, test matchTest, node ast.Node) {
	pat, err := Compile(test.tmpl)
	if err != nil {
		t.Errorf("compile(%q): %v", test.tmpl, err)
		return
	}

	if shouldAccept {
		if !pat.Match(node) {
			t.Errorf("%q.accept(%q) failed", test.tmpl, test.input)
		}
	} else {
		if pat.Match(node) {
			t.Errorf("%q.reject(%q) failed", test.tmpl, test.input)
		}
	}
}

var astNodes = make(map[string]ast.Node)

func astStmt(t *testing.T, code string) ast.Node {
	if node := astNodes[code]; node != nil {
		return node
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(
		fset, "", "package main;func main() {"+code+"}", parser.ParseComments)
	if err != nil {
		t.Fatalf("Go parse stmt: %v", err)
	}
	node := f.Decls[0].(*ast.FuncDecl).Body.List[0]
	astNodes[code] = node
	return node
}

func astExpr(t *testing.T, code string) ast.Node {
	if node := astNodes[code]; node != nil {
		return node
	}
	node, err := parser.ParseExpr(code)
	if err != nil {
		t.Fatalf("Go parse expr: %v", err)
	}
	astNodes[code] = node
	return node
}

func astFile(t *testing.T, code string) ast.Node {
	if node := astNodes[code]; node != nil {
		return node
	}
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Go parse file: %v", err)
	}
	astNodes[code] = node
	return node
}

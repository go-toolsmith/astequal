package astcmp

import (
	"testing"
)

type astEqualTest struct {
	// x/y - expr/stmt/decl strings.

	x string
	y string

	want bool
}

func TestEqualExprNils(t *testing.T) {
	t.SkipNow() // TODO
}

func TestEqualStmtNils(t *testing.T) {
	t.SkipNow() // TODO
}

func TestEqualDeclNils(t *testing.T) {
	t.SkipNow() // TODO
}

func TestEqualExprString(t *testing.T) {
	runTest := func(t *testing.T, test astEqualTest) {
		have := EqualExprString(test.x, test.y)
		if have != test.want {
			t.Errorf("EqualExprString:\nx: %q\ny: %q\nhave: %v\nwant: %v",
				test.x, test.y, have, test.want)
			return
		}
		have = Equal(astParseExpr(test.x), astParseExpr(test.y))
		if have != test.want {
			t.Errorf("Equal:\nx: %q\ny: %q\nhave: %v\nwant: %v",
				test.x, test.y, have, test.want)
		}
	}
	runTests := func(name string, tests []astEqualTest) {
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				runTest(t, test)
			}
		})
	}
	test := func(name string, expressions ...string) {
		var suite []astEqualTest

		for _, x := range expressions {
			suite = append(suite, astEqualTest{
				x:    x,
				y:    `/**/` + x + `/**/`,
				want: true,
			})
			suite = append(suite, astEqualTest{x: x, y: `a[` + x + `]`})
			for _, y := range expressions {
				if x == y {
					continue
				}
				test := astEqualTest{x: x + "/**/", y: "/**/" + y}
				suite = append(suite, test)
			}
		}

		runTests(name, suite)
	}

	runTests("Malformed", []astEqualTest{
		// Malformed expressions are replaced with BadExpr,
		// which never compared as equal.
		{`------`, `------`, false},
		{`?!`, `?!`, false},
		{`$#`, `$#`, false},
		{`x===y`, `x===y`, false},
		{`@a`, `@a`, false},
	})

	test("Ident", `x`, `y`, `abc`)

	test("BasicLit",
		`1`, `250`,
		`1.2`, `0.77`,
		`1+2i`, `0.i`,
		`'a'`, `'ä'`,
		`"abc"`, `"0"`)

	test("FuncLit",
		`func() {}`,
		`func(x int) {}`,
		`func(y int) {}`,
		`func(x ...int) {}`,
		`func(x int) int {}`,
		`func(x int) float64 {}`,
		`func(x, y int) {}`,
		`func(y, x int) {}`,
		`func(x int) (a int, b int) {}`,
		`func(x int) (int, rune) {}`,
		`func(x int) (rune, int) {}`)

	test("CompositeLit",
		`X{}`,
		`Y{}`,
		`X{a: 1}`,
		`X{b: 1}`,
		`X{1, 2}`,
		`X{2, 1}`,
		`X{a: A{1}}`,
		`X{a: B{1}}`,
		`struct{}{}`,
		`struct{a int}{1}`,
		`struct{a, b int}{1, 2}`,
		`[...]int{1}`,
		`[1]int{1}`,
		`[]int{1}`)

	test("ParenExpr", `(1)`, `((1))`, `(((x)))`)

	test("SelectorExpr",
		`a.b`,
		`a.x`,
		`a.b.c`,
		`a.b.x`,
		`a().b`,
		`a().b()`)

	test("IndexExpr",
		`a[0]`,
		`b[0]`,
		`a[1]`,
		`a[x][y]`,
		`a[y][x]`,
		`a[a[x]]`,
		`a[b[x]]`)

	test("SliceExpr",
		`a[x:y]`,
		`b[x:y]`,
		`a[y:x]`,
		`a[x:y:a]`,
		`a[x:y:b]`,
		`a[x:]`,
		`a[:x]`,
		`a[:]`,
		`a[x][1:]`,
		`a[x][2:]`)

	test("TypeAssertExpr",
		`x.(int)`,
		`y.(int)`,
		`x.([]int)`,
		`x.(float32)`,
		`x.(a).(*b)`,
		`x.(*[2]interface{})`,
		`x.([2]interface{})`)

	test("CallExpr",
		`int(x)`,
		`int(y)`,
		`f(x)`,
		`(int)(x)`,
		`(uintptr)(unsafe.Pointer(y))`,
		`f(a, xs)`,
		`f(a, xs...)`)

	test("StarExpr", `*x`, `*y`, `**x`)

	test("ArrayType",
		`[10][x]int`,
		`[x+y]int`,
		`[x+y]float32`,
		`[][]int`,
		`[]int`,
		`[]float32`)

	test("StructType",
		`struct{}`,
		`struct{x int}`,
		`struct{y int}`,
		`struct{x int; y string}`,
		`struct{x, y int}`)

	test("FuncType",
		`func()`,
		`func(A)`,
		`func(B)`,
		`func(A) A`,
		`func(A) B`,
		`func(A) (A, A)`,
		`func(A, ...B)`,
		`func(A, B)`)

	test("InterfaceType",
		`interface{}`,
		`interface{A()}`,
		`interface{B()}`,
		`interface{A(); B()}`)

	test("MapType",
		`map[A]B`,
		`map[B]A`,
		`map[A]map[B]C`,
		`map[A][]map[B]C`)

	test("ChanType",
		`chan A`,
		`chan B`,
		`<- chan A`,
		`chan <- A`,
		`chan chan <- A`,
		`chan <- chan A`)
}

func TestEqualStmtString(t *testing.T) {
	runTest := func(t *testing.T, test astEqualTest) {
		have := EqualStmtString(test.x, test.y)
		if have != test.want {
			t.Errorf("EqualStmtString:\nx: %q\ny: %q\nhave: %v\nwant: %v",
				test.x, test.y, have, test.want)
			return
		}
		have = Equal(astParseStmt(test.x), astParseStmt(test.y))
		if have != test.want {
			t.Errorf("Equal:\nx: %q\ny: %q\nhave: %v\nwant: %v",
				test.x, test.y, have, test.want)
		}
	}
	runTests := func(name string, tests []astEqualTest) {
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				runTest(t, test)
			}
		})
	}
	test := func(name string, expressions ...string) {
		var suite []astEqualTest

		for _, x := range expressions {
			suite = append(suite, astEqualTest{
				x:    x,
				y:    `/**/` + x + `/**/`,
				want: true,
			})
			suite = append(suite, astEqualTest{x: x, y: `{` + x + `}`})
			for _, y := range expressions {
				if x == y {
					continue
				}
				test := astEqualTest{x: x + "/**/", y: "/**/" + y}
				suite = append(suite, test)
			}
		}

		runTests(name, suite)
	}

	runTests("Malformed", []astEqualTest{
		// Malformed statements are replaced with BadStmt,
		// which never compared as equal.
		{`------`, `------`, false},
		{`?!`, `?!`, false},
		{`$#`, `$#`, false},
		{`x===y`, `x===y`, false},
		{`@a`, `@a`, false},
		{`{------}`, `{------}`, false},
		{`{?!}`, `{?!}`, false},
		{`{$#}`, `{$#}`, false},
		{`{x===y}`, `{x===y}`, false},
		{`{@a}`, `{@a}`, false},
	})

	test("EmptyStmt", `;`)

	test("ExprStmt",
		`f()`,
		`x.y()`,
		`<- x`,
		`<- y`,
		`<- <- x`)

	test("IncDecStmt",
		`x++`,
		`y++`,
		`x--`,
		`y--`)

	test("AssignStmt",
		`x := y`,
		`x = y`,
		`y := x`,
		`y = x`,
		`x, y := a, b`,
		`x, y = a, b`,
		`x, y := b, a`,
		`x, y = b, a`,
		`x, y := A()`,
		`x, y, z := A()`)

	test("GoStmt",
		`go A()`,
		`go B()`,
		`go func(){}()`,
		`go f()(x)`,
		`go f()(y)`)

	test("DeferStmt",
		`defer f()`,
		`defer g()`,
		`defer func(){}()`)

	test("ReturnStmt",
		`return`,
		`return x`,
		`return x, y`,
		`return y, x`)

	test("SwitchStmt",
		`switch {}`,
		`switch {default: f()}`,
		`switch {default: f(); break}`,
		`switch {case x == y: f()}`,
		`switch x {case A: f()}`,
		`switch y {case A: f()}`,
		`switch x {case B: f()}`,
		`switch y := x; y {case 0: f()}`)

	test("TypeSwitchStmt",
		`switch x := a.(type) {}`,
		`switch x := b.(type) {}`,
		`switch y := a.(type) {}`,
		`switch x := a.(type) {case int: f()}`,
		`switch x := a.(type) {default: f()}`,
		`switch a.(type) {}`,
		`switch b.(type) {}`,
		`switch a.(type) {case int: f()}`,
		`switch a.(type) {default: f()}`)

	test("IfStmt",
		`if x {}`,
		`if x {f()}`,
		`if y {}`,
		`if y {f()}`,
		`if x {f()} else {g()}`,
		`if x {f()} else {f()}`,
		`if y := x; y {f()}`,
		`if x := y; x {f()}`,
		`if x {f()} else if y {g()} else {}`,
		`if x {f()} else if y {g()} else {panic(0)}`)

	test("SelectStmt",
		`select {}`,
		`select {case x <- a: f()}`,
		`select {case x <- a: f(); case <-b: g()}`,
		`select {case <- a: f()}`,
		`select {case <- b: f()}`,
		`select {case <- a: g()}`)

	test("ForStmt",
		`for {}`,
		`for true {}`,
		`for i := 0; i < x; i++ {}`,
		`for i := 0; i < len(xs); i++ {}`,
		`for i, j := f(); i > 0; i, j = x+i, y+j {}`,
		`for i = 0; i < x; i++ {}`,
		`for i = 0; i < len(xs); i++ {}`,
		`for i, j = f(); i > 0; i, j = x+i, y+j {}`)

	test("RangeStmt",
		`for range xs {}`,
		`for range ys {}`,
		`for a, b := range xs {}`,
		`for b, a := range xs {}`,
		`for _, b := range xs {}`,
		`for a, _ := range xs {}`,
		`for a := range xs {}`,
		`for a, b = range xs {}`,
		`for b, a = range xs {}`,
		`for _, b = range xs {}`,
		`for a, _ = range xs {}`,
		`for a = range xs {}`)

	test("BranchStmt",
		`fallthrough`,
		`goto x`,
		`goto y`,
		`break`,
		`break x`,
		`break y`,
		`continue`,
		`continue x`,
		`continue y`)

	test("DeclStmt",
		`var x A`,
		`var x B`,
		`var y A`,
		`var x, y A`,
		`var y, x A`,
		`var (x A; y A; z A)`,
		`var (x, y, z A)`,
		`var x A = a`,
		`var x A = b`,
		`const x = a`,
		`var x = a`)

	test("LabeledStmt",
		`x: {}`,
		`x: for {}`,
		`y: {}`,
		`y: for {}`)
}

func TestEqualDeclString(t *testing.T) {
	runTest := func(t *testing.T, test astEqualTest) {
		have := EqualDeclString(test.x, test.y)
		if have != test.want {
			t.Errorf("EqualDeclString:\nx: %q\ny: %q\nhave: %v\nwant: %v",
				test.x, test.y, have, test.want)
			return
		}
		have = Equal(astParseDecl(test.x), astParseDecl(test.y))
		if have != test.want {
			t.Errorf("Equal:\nx: %q\ny: %q\nhave: %v\nwant: %v",
				test.x, test.y, have, test.want)
		}
	}
	runTests := func(name string, tests []astEqualTest) {
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				runTest(t, test)
			}
		})
	}
	test := func(name string, expressions ...string) {
		var suite []astEqualTest

		for _, x := range expressions {
			suite = append(suite, astEqualTest{
				x:    x,
				y:    `/**/` + x + `/**/`,
				want: true,
			})
			for _, y := range expressions {
				if x == y {
					continue
				}
				test := astEqualTest{x: x + "/**/", y: "/**/" + y}
				suite = append(suite, test)
			}
		}

		runTests(name, suite)
	}

	runTests("Malformed", []astEqualTest{
		// Malformed declarations are replaced with BadDecl,
		// which never compared as equal.
		{`------`, `------`, false},
		{`?!`, `?!`, false},
		{`$#`, `$#`, false},
		{`x===y`, `x===y`, false},
		{`@a`, `@a`, false},
	})

	test("ImportSpec",
		`import ()`,
		`import "a"`,
		`import "b"`,
		`import (. "a")`,
		`import ("a"; "b")`,
		`import (x "a"; y "b")`,
		`import (y "a"; x "b")`,
		`import ("b"; "a")`,
		`import ("a"; "b"; "c")`)

	test("TypeSpec",
		`type ()`,
		`type x A`,
		`type y A`,
		`type x B`,
		`type x struct {a int; b int}`,
		`type x struct {b int; a int}`,
		`type x interface {f()}`,
		`type x interface {g()}`,
		`type (x A; y B)`,
		`type (x A; y B; z C)`)

	test("ValueSpec",
		`var x A`,
		`var x B`,
		`var y A`,
		`var x, y A`,
		`var y, x A`,
		`var (x A; y A; z A)`,
		`var (x, y, z A)`,
		`var x A = a`,
		`var x A = b`,
		`const x = a`,
		`var x = a`)

	test("FuncDecl",
		`func f() {}`,
		`func g() {}`,
		`func f(x, y int) int {return 0}`,
		`func f(y, x int) int {return 0}`,
		`func (A) f() {}`,
		`func (B) f() {}`,
		`func (a A) f() {}`,
		`func (b B) f() {}`)
}
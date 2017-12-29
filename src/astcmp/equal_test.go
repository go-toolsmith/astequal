package astcmp

import (
	"testing"
)

type astEqualTest struct {
	// X/Y - expr/stmt/decl strings
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
	yes := func(x string) astEqualTest {
		return astEqualTest{x, "\n/**/" + x + "  \n//xyz\n", true}
	}
	no := func(x, y string) astEqualTest { return astEqualTest{x, y, false} }

	tests := []astEqualTest{
		// Malformed expressions are replaced with BadExpr,
		// which never compared as equal.
		{`------`, `------`, false},
		{`?!`, `?!`, false},
		{`$#`, `$#`, false},
		{`x===y`, `x===y`, false},
		{`@a`, `@a`, false},

		yes(`x`),
		yes(`foo`),
		no(`x`, `y`),
		no(`foo`, `y`),
		no(`x`, `(x)`),

		yes(`1`),
		yes(`250`),
		no(`1`, `2`),
		no(`1`, `777`),
		no(`1`, `(1)`),

		yes(`1.2`),
		yes(`250.8900`),
		no(`1`, `1.0`),
		no(`1.01`, `1.00`),
		no(`1.0`, `(1.0)`),

		yes(`1+2i`),
		yes(`0i`),
		yes(`0.i`),
		yes(`.12345E+5i`),
		no(`1+2i `, `1+3i`),
		no(`0i`, `1i`),
		no(`0.i`, `1.i`),
		no(`.12345E+5i`, `.12245E+5i`),
		no(`0i`, `(0i)`),

		yes(`'a'`),
		yes(`'\''`),
		yes(`'本'`),
		yes(`'\377'`),
		yes(`'\U00100230'`),
		no(`'a'`, `'b'`),
		no(`'a'`, `'ä'`),
		no(`'\377'`, `'\370'`),
		no(`'\U00101234'`, `'\U00100230'`),
		no(`'a'`, `('a')`),

		yes(`func(){}`),
		yes(`func(x int){}`),
		yes(`func(x, y int){}`),
		yes(`func(x, y int)float32{}`),
		yes(`func(x, y int) float32 {println(1)}`),
		yes(`func(x int, y string) (int, rune) {return 0, '0'}`),
		no(`func(){}`, `func(y int){}`),
		no(`func(x int){}`, `func(y int){}`),
		no(`func(x, y int){}`, `func(x int, y int){}`),
		no(`func(x, y int)float32{}`, `func(x, y int)string{}`),
		no(
			`func(x, y int) float32 {println(1)}`,
			`func(x, y int) float32 {println(0)}`,
		),
		no(
			`func(x int, y string) (int, rune) {return 0, '0'}`,
			`func(x int, y string) (int, int) {return 0, '0'}`,
		),
		no(`func(){}`, `(func(){})`),

		yes(`Foo{}`),
		yes(`Foo{x:1}`),
		yes(`Foo{1,2}`),
		yes(`struct{}{}`),
		yes(`struct{x int}{1}`),
		yes(`struct{x, y int}{1, 2}`),
		no(`Foo{}`, `Bar{}`),
		no(`Foo{x:1}`, `Foo{x:2}`),
		no(`Foo{x:1}`, `Foo{y:2}`),
		no(`Foo{1,2}`, `Foo{2,2}`),
		no(`struct{}{}`, `struct{x int}{}`),
		no(`struct{x int}{1}`, `struct{x int}{2}`),
		no(`struct{x, y int}{1, 2}`, `struct{x, z int}{1, 2}`),
		no(`Foo{}`, `(Foo{})`),

		yes(`(1)`),
		yes(`((1))`),
		yes(`(((a)))`),
		no(`(1)`, `1`),
		no(`((1))`, `(1)`),
		no(`(((a)))`, `((a))`),

		yes(`a.b`),
		yes(`a.b.c`),
		yes(`foo().bar`),
		yes(`foo().bar()`),
		no(`a.b`, `a.c`),
		no(`a.b.c`, `c.b.c`),
		no(`foo().bar`, `foo().baz`),
		no(`foo().bar()`, `foo().baz()`),
		no(`a.b`, `(a.b)`),

		yes(`a[0]`),
		yes(`a[0][x]`),
		yes(`a[0][f()]`),
		yes(`a[0].x`),
		no(`a[0]`, `b[0]`),
		no(`a[0][x]`, `a[0][y]`),
		no(`a[0][f()]`, `a[0][g()]`),
		no(`a[0].x`, `a[1].x`),
		no(`a[0]`, `(a[0])`),

		yes(`a[x:y]`),
		yes(`a[x:y:z]`),
		yes(`a[x:][:]`),
		yes(`a[:x][:]`),
		yes(`a[:][:]`),
		no(`b[x:y]`, `a[x:y]`),
		no(`b[x:y:z]`, `a[x:y:z]`),
		no(`b[x:][:]`, `a[x:][:]`),
		no(`b[:x][:]`, `a[:x][:]`),
		no(`b[:][:]`, `a[:][:]`),
		no(`a[x:y]`, `a[x:x]`),
		no(`a[x:y:z]`, `a[x:y:x]`),
		no(`a[x:][:]`, `a[x:][:1]`),
		no(`a[:x][:]`, `a[:x][1:]`),
		no(`a[:][:]`, `a[1:][:]`),
		no(`a[x:y]`, `(a[x:y])`),

		yes(`x.(int)`),
		yes(`x.(a).(*b)`),
		yes(`x.([]int)`),
		yes(`x.(*[2]interface{})`),
		no(`x.(int)`, `x.(float)`),
		no(`x.(int)`, `y.(int)`),
		no(`x.(a).(*b)`, `x.(a).(b)`),
		no(`x.([]int)`, `x.([10]int)`),
		no(`x.(*[2]interface{})`, `x.(*[2]struct{})`),
		no(`x.(int)`, `(x.(int))`),

		yes(`int(x)`),
		yes(`(*int)(x)`),
		yes(`(uintptr)(unsafe.Pointer(y))`),
		no(`int(x)`, `float32(y)`),
		no(`(*int)(x)`, `(int)(x)`),
		no(`(uintptr)(unsafe.Pointer(y))`, `uintptr(unsafe.Pointer(y))`),
		no(`(*int)(x)`, `((*int)(x))`),

		yes(`f()`),
		yes(`f(1,2)`),
		yes(`f(1,xs...)`),
		yes(`f(g())`),
		no(`f()`, `g()`),
		no(`f(1,2)`, `f(1,3)`),
		no(`f(1,xs...)`, `f(1,xs)`),
		no(`f(g())`, `f(f())`),

		yes(`*x`),
		yes(`**x`),
		no(`*x`, `*y`),
		no(`**x`, `*x`),
		no(`*x`, `(*x)`),

		yes(`+x`),
		yes(`-x`),
		no(`+x`, `+y`),
		no(`+x`, `-x`),
		no(`+x`, `(+x)`),

		yes(`x+y`),
		yes(`x-y`),
		yes(`x+y+z`),
		yes(`x+y-z`),
		no(`x+y`, `0+y`),
		no(`x-y`, `x-0`),
		no(`x+y+z`, `x+y+0`),
		no(`x+y-z`, `x+y-0`),
		no(`x+y`, `(x+y)`),

		yes(`[10][x]int`),
		yes(`[x+y]int`),
		yes(`[]int`),
		yes(`[][]int`),
		no(`[10][x]int`, `[10][y]int`),
		no(`[x+y]int`, `[x+y]float32`),
		no(`[]int`, `[][]int`),
		no(`[][]int`, `[]int`),
		no(`[]int`, `([]int)`),

		yes(`[...]int{1}`),
		yes(`struct{x int}`),
		yes(`struct{}`),
		yes(`struct{x int; y string}`),
		no(`[...]int{1}`, `[1]int{1}`),
		no(`struct{x int}`, `struct{y int}`),
		no(`struct{}`, `struct{x int}`),
		no(`struct{x int; y string}`, `struct{x float32; y string}`),
		no(`struct{}`, `(struct{})`),

		yes(`func()`),
		yes(`func(int)`),
		yes(`func(int)(int, int)`),
		yes(`func(int)(int, string)`),
		yes(`func(float32, float64)int`),
		no(`func()`, `func(int)`),
		no(`func(int)`, `func(int,int)`),
		no(`func(int)(int, int)`, `func(int)(float32, int)`),
		no(`func(int)(int, int)`, `func(int)(int)`),
		no(`func(float32, float64)int`, `func(float64, float64)`),
		no(`func()`, `(func())`),

		yes(`interface{foo()}`),
		yes(`interface{foo(); bar(int)}`),
		yes(`interface{}`),
		no(`interface{foo()}`, `interface{bar()}`),
		no(`interface{foo(); bar(int)}`, `interface{foo(); bar(float32)}`),
		no(`interface{}`, `interface{foo(int)}`),
		no(`interface{}`, `(interface{})`),

		yes(`map[A]B`),
		yes(`map[A]map[int]int`),
		yes(`map[A][]int`),
		no(`map[A]B`, `map[A]int`),
		no(`map[A]B`, `map[int]B`),
		no(`map[A]map[int]int`, `map[A]map[int]float32`),
		no(`map[A][]int`, `map[A][2]int`),
		no(`map[A][B]`, `(map[A][B])`),

		yes(`chan int`),
		yes(`<- chan int`),
		yes(`chan <- int`),
		yes(`chan chan <- int`),
		no(`chan int`, `chan x`),
		no(`chan <- int`, `<- chan int`),
		no(`chan <- x`, `chan <- int`),
		no(`chan chan <- int`, `chan <- int`),
		no(`chan int`, `(chan int)`),
	}

	for _, test := range tests {
		have := EqualExprString(test.x, test.y)
		if have != test.want {
			t.Errorf("EqualExprString:\nx: %q\ny: %q\nhave: %v\nwant: %v",
				test.x, test.y, have, test.want)
		}
		x, y := astParseExpr(test.x), astParseExpr(test.y)
		have = Equal(x, y)
		if have != test.want {
			t.Errorf("Equal:\nx: %q\ny: %q\nhave: %v\nwant: %v",
				test.x, test.y, have, test.want)
		}
	}
}

func TestEqualStmtString(t *testing.T) {
	yes := func(x string) astEqualTest {
		return astEqualTest{"{" + x + "}", "\n/**/{" + x + " } \n//xyz\n", true}
	}
	no := func(x, y string) astEqualTest {
		return astEqualTest{"{" + x + "}", "{" + y + "}", false}
	}

	tests := []astEqualTest{
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

		yes(`;`),
		yes(`;;`),
		no(`{;}`, `{}`),
		no(`{f(), g()}`, `{f()}`),

		yes(`f()`),
		yes(`x.y()`),
		no(`f()`, `g()`),
		no(`x.y()`, `x.z()`),
		no(`f()`, `{f()}`),

		yes(`<- x`),
		yes(`<- <- x`),
		no(`<- x`, `<- y`),
		no(`<- <- x`, `<- x`),
		no(`<- x`, `{<- x}`),

		yes(`x++`),
		yes(`x--`),
		no(`x++`, `y++`),
		no(`x--`, `x++`),
		no(`x++`, `{x++}`),

		yes(`x := y`),
		yes(`x = y`),
		yes(`x, y := 1, 2`),
		yes(`x, y = 1, 2`),
		yes(`x, y := f()`),
		yes(`x, y = f()`),
		no(`x := y`, `x := z`),
		no(`x = y`, `x := y`),
		no(`x, y := 1, 2`, `x, y = 2, 1`),
		no(`x, y = 1, 2`, `x = 1`),
		no(`x, y := f()`, `x, y := g()`),
		no(`x, y = f()`, `x := f()`),
		no(`x := y`, `{x := y}`),

		yes(`go f()`),
		yes(`go func(){}()`),
		yes(`go f()(1)`),
		no(`go f()`, `go g()`),
		no(`go func(){}()`, `go func()int{}()`),
		no(`go f()(1)`, `go f(1)()`),
		no(`go f()`, `{go f()}`),

		yes(`defer f()`),
		no(`defer f()`, `defer g()`),
		no(`defer f()`, `{defer f()}`),

		yes(`return`),
		yes(`return x`),
		yes(`return x, y`),
		no(`return`, `return x`),
		no(`return x`, `return y`),
		no(`return x, y`, `return x`),
		no(`return x, y`, `return y, x`),
		no(`return x`, `{return x}`),

		yes(`switch x {case 0: f()}`),
		yes(`switch y := x; y {case 0: f()}`),
		yes(`switch {case x == 0: f()}`),
		yes(`switch {default: f()}`),
		yes(`switch {default: break}`),
		no(
			`switch x {case 0: f()}`,
			`switch y {case 0: f()}`,
		),
		no(
			`switch x {case 0: f()}`,
			`switch x {case 1: f()}`,
		),
		no(
			`switch x {case 0: f()}`,
			`switch x {case 0: g()}`,
		),
		no(
			`switch y := x; y {case 0: f()}`,
			`switch y := z; y {case 0: f()}`,
		),
		no(
			`switch y := x; y {case 0: f()}`,
			`switch z := x; z {case 0: f()}`,
		),
		no(
			`switch y := x; y {case 0: f()}`,
			`switch y := x; y {case 1: f()}`,
		),
		no(
			`switch y := x; y {case 0: f()}`,
			`switch y := x; y {case 0: g()}`,
		),
		no(
			`switch {case x == 0: f()}`,
			`switch {case x == 1: f()}`,
		),
		no(
			`switch {case x == 0: f()}`,
			`switch x {case 0: f()}`,
		),
		no(
			`switch {default: break}`,
			`switch {default: f()}`,
		),
		no(`switch{}`, `{switch{}}`),

		yes(`if true {f()}`),
		yes(`if x := y; x {f()}`),
		yes(`if x {f()} else {g()}`),
		yes(`if x {f()} else if y {g()} else {panic(0)}`),
		no(`if true {f()}`, `if false {f()}`),
		no(`if true {f()}`, `if true {g()}`),
		no(`if x := y; x {f()}`, `if x := z; x {f()}`),
		no(`if x {f()} else {g()}`, `if x {f()} else {f()}`),
		no(`if x {f()} else {g()}`, `if x {f()}`),
		no(
			`if x {f()} else if y {g()} else {panic(0)}`,
			`if x {f()} else if y {g()} else {panic(1)}`,
		),
		no(`if true {f()}`, `{if true {f()}}`),

		yes(`switch x := y.(type) {}`),
		yes(`switch x.(type) {}`),
		yes(`switch x := y.(type) {case int: f()}`),
		yes(`switch x.(type) {default: g()}`),
		no(`switch x := y.(type) {}`, `switch x := z.(type) {}`),
		no(`switch x.(type) {}`, `switch y.(type) {}`),
		no(`switch x.(type) {}`, `{switch x.(type) {}}`),

		yes(`select {}`),
		yes(`select {case x <- c: f(); case <-y: g()}`),
		no(
			`select {case x <- c: f(); case <-y: g()}`,
			`select {case <-x: f(); case <-y: g()}`,
		),
		no(`select {}`, `{select {}}`),

		yes(`for {}`),
		yes(`for true {}`),
		yes(`for ;; {}`),
		yes(`for i := 0; i < 10; i++ {}`),
		yes(`for i = 0; i < len(xs); i += 2 {}`),
		no(`for true {}`, `for false {}`),
		no(`for ;; {}`, `for ; i < 10; {}`),
		no(`for ;; {}`, `for true {}`),
		no(
			`for i := 0; i < 10; i++ {}`,
			`for i := 0; i < 10; j++ {}`,
		),
		no(
			`for i := 0; i < 10; i++ {}`,
			`for i := 0; j < 10; i++ {}`,
		),
		no(
			`for i := 0; i < 10; i++ {}`,
			`for j := 0; i < 10; i++ {}`,
		),
		no(
			`for i = 0; i < len(xs); i += 2 {}`,
			`for i := 0; i < len(xs); i += 2 {}`,
		),
		no(`for {}`, `{for {}}`),

		yes(`for i := range m {}`),
		yes(`for k, v := range m {}`),
		yes(`for _, v := range m {}`),
		yes(`for k, _ := range m {}`),
		yes(`for i = range m {}`),
		yes(`for k, v = range m {}`),
		no(`for i := range m {}`, `for i := range mm {}`),
		no(`for k, v := range m {}`, `for _, v := range m {}`),
		no(`for _, v := range m {}`, `for k, _ := range m {}`),
		no(`for k, _ := range m {}`, `for _, _ := range m {}`),
		no(`for i = range m {f()}`, `for i := range m {g()}`),
		no(`for k, v = range m {}`, `for k, v := range m {}`),
		no(`for i := range m {}`, `{for i := range m{}}`),

		yes(`goto x`),
		yes(`break x`),
		no(`goto x`, `{goto x}`),

		yes(`var x int`),
		yes(`var x = 10`),
		yes(`var (x int; y int)`),
		yes(`const x = 10`),
		no(`var x = 10`, `var y = 10`),
		no(`var x = f()`, `var x, y = f()`),
		no(`const x = 10`, `var x = 10`),
		no(`var x int`, `{var x int}`),

		yes(`l: for {}`),
		no(`a: for{}`, `b: for{}`),
	}

	for _, test := range tests {
		have := EqualStmtString(test.x, test.y)
		if have != test.want {
			t.Errorf("EqualStmtString:\nx: %q\ny: %q\nhave: %v\nwant: %v",
				test.x, test.y, have, test.want)
		}
		x, y := astParseStmt(test.x), astParseStmt(test.y)
		have = Equal(x, y)
		if have != test.want {
			t.Errorf("Equal:\nx: %q\ny: %q\nhave: %v\nwant: %v",
				test.x, test.y, have, test.want)
		}
	}
}

func TestEqualDeclString(t *testing.T) {
	yes := func(x string) astEqualTest {
		return astEqualTest{x, "\n" + x + "\n", true}
	}
	no := func(x, y string) astEqualTest {
		return astEqualTest{x, y, false}
	}

	tests := []astEqualTest{
		// Malformed declarations are replaced with BadDecl,
		// which never compared as equal.
		{`------`, `------`, false},
		{`?!`, `?!`, false},
		{`$#`, `$#`, false},
		{`x===y`, `x===y`, false},
		{`@a`, `@a`, false},

		yes(`import ()`),
		yes(`import "a"`),
		yes(`import ("a"; "b")`),
		no(`import "a"`, `import "b"`),
		no(`import ("a"; "b")`, `import ("a", "a")`),
		no(`import ("a"; "b")`, `import ("a")`),
		no(`import ()`, `func f() {}`),

		yes(`func f() {}`),
		yes(`func f(x, y int) int {return 0}`),
		yes(`func (Foo) f() {}`),
		yes(`func (x Foo) f() {}`),
		no(`func f() {}`, `func g() {}`),
		no(`func f() {}`, `func f() {return}`),
		no(`func (Foo) f() {}`, `func (x Foo) f() {}`),
		no(`func (Foo) f() {}`, `func (Bar) f() {}`),
		no(`func (x Foo) f() {}`, `func (y Foo) f() {}`),

		yes(`type x int`),
		yes(`type x struct{foo int}`),
		yes(`type x interface{foo()}`),
		no(`type x int`, `type y int`),
		no(`type x int`, `type x float32`),
		no(`type x interface{foo()}`, `type x interface{bar()}`),
		no(`type x interface{}`, `type x struct{}`),
	}

	for _, test := range tests {
		have := EqualDeclString(test.x, test.y)
		if have != test.want {
			t.Errorf("EqualDeclString:\nx: %q\ny: %q\nhave: %v\nwant: %v",
				test.x, test.y, have, test.want)
		}
		x, y := astParseDecl(test.x), astParseDecl(test.y)
		have = Equal(x, y)
		if have != test.want {
			t.Errorf("Equal:\nx: %q\ny: %q\nhave: %v\nwant: %v",
				test.x, test.y, have, test.want)
		}
	}
}

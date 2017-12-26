package astexp

import (
	"testing"
)

func BenchmarkCompile(b *testing.B) {
	templates := []string{
		`(* $int $int)`,
		`(block (call f) (call g 1) (+ (+ (+ 1 3) 5) 6))`,
		`(array_type 10 (array_type 20 (slice_type int)))`,
		`(call f $int $... $str $... x y z 0 1 2)`,
	}
	cl := NewCompiler()
	for i := 0; i < b.N; i++ {
		for _, tmpl := range templates {
			_, err := cl.CompileString(tmpl)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

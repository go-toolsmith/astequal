package astexp

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"text/scanner"
	"unicode"
)

// Compiler constructs Pattern from string sexp template.
//
// Should be created with NewCompiler function.
type Compiler struct {
	scanner.Scanner
}

// NewCompiler returns compiler for sexp.
func NewCompiler() *Compiler {
	var cl Compiler

	// Extend default ident definition with '$' and '.'.
	cl.IsIdentRune = func(ch rune, i int) bool {
		return unicode.IsLetter(ch) ||
			(i != 0 && unicode.IsDigit(ch)) ||
			(i != 0 && ch == '.') ||
			ch == '_' ||
			ch == '$'
	}

	return &cl
}

// CompileString builds Pattern from given sexp template.
// Convenience wrapper over Compiler.Compile.
func (cl *Compiler) CompileString(sexp string) (*Pattern, error) {
	return cl.Compile(strings.NewReader(sexp))
}

// Compile builds Pattern from sexp template that can be read from r.
func (cl *Compiler) Compile(r io.Reader) (*Pattern, error) {
	cl.Init(r)
	return cl.compile()
}

func (cl *Compiler) compile() (*Pattern, error) {
	pat, err := cl.compileTok(cl.Scan())
	if err != nil {
		return nil, err
	}
	// This check is not strictly necessary, but can help users to detect
	// unintuitive behavior for templates like "1 2 3" which will match only
	// first atom.
	if cl.Peek() != scanner.EOF {
		return nil, errors.New("unconsumed input after sexp template (forgot to forms into list?)")
	}
	return &pat, nil
}

func (cl *Compiler) compileWildcard(name, kind string) (Pattern, error) {
	if name != "" {
		return Pattern{}, errors.New("named wildcards are not implemented")
	}
	switch kind {
	case "...":
		return Pattern{tag: tagEllipsis}, nil

	case "id":
		return Pattern{tag: tagAnyIdent}, nil
	case "int":
		return Pattern{tag: tagAnyInt}, nil
	case "float":
		return Pattern{tag: tagAnyFloat}, nil
	case "imag":
		return Pattern{tag: tagAnyImag}, nil
	case "char":
		return Pattern{tag: tagAnyChar}, nil
	case "str":
		return Pattern{tag: tagAnyString}, nil

	default:
		return Pattern{}, fmt.Errorf("unknown wildcard: %v", kind)
	}
}

func (cl *Compiler) compileIdent() (Pattern, error) {
	name := cl.TokenText()
	if name[0] == '$' {
		name = name[len("$"):]
		if cl.tryConsume(':') {
			if cl.Scan() != scanner.Ident {
				return Pattern{}, errors.New("wildcard: expected ident")
			}
			return cl.compileWildcard(name, cl.TokenText())
		}
		return cl.compileWildcard("", name)
	}
	return Pattern{tag: tagIdent, value: name}, nil
}

func (cl *Compiler) compileTok(tok rune) (Pattern, error) {
	switch tok {
	case '(':
		return cl.compileList()
	case scanner.Ident:
		return cl.compileIdent()

	case scanner.Int:
		return Pattern{tag: tagInt, value: cl.TokenText()}, nil
	case scanner.Float:
		return Pattern{tag: tagFloat, value: cl.TokenText()}, nil
	case scanner.Char:
		return Pattern{tag: tagChar, value: cl.TokenText()}, nil
	case scanner.String:
		return Pattern{tag: tagString, value: cl.TokenText()}, nil

	case scanner.EOF:
		return Pattern{}, errors.New("unexpected token: EOF")
	default:
		return Pattern{}, fmt.Errorf("unexpected token: %q", tok)
	}
}

func (cl *Compiler) compileList() (Pattern, error) {
	tag, err := cl.compileListTag()
	if err != nil {
		return Pattern{}, fmt.Errorf("list head: %v", err)
	}

	pat := Pattern{tag: tag}
	for tok := cl.Scan(); tok != ')'; tok = cl.Scan() {
		if tok == scanner.EOF {
			return Pattern{}, errors.New("unterminated list")
		}
		x, err := cl.compileTok(tok)
		if err != nil {
			return Pattern{}, fmt.Errorf("list: %v", err)
		}
		pat.tail = append(pat.tail, x)
	}
	return pat, nil
}

func (cl *Compiler) tryConsume(tok rune) bool {
	if cl.Peek() == tok {
		cl.Scan()
		return true
	}
	return false
}

func (cl *Compiler) compileListTag() (patternTag, error) {
	switch cl.Scan() {
	case '+':
		if cl.tryConsume('+') {
			return tagOpInc, nil
		}
		return tagOpAdd, nil
	case '-':
		if cl.tryConsume('-') {
			return tagOpDec, nil
		}
		return tagOpSub, nil
	case '*':
		return tagOpMul, nil
	case '/':
		return tagOpQuo, nil
	case '%':
		return tagOpRem, nil
	case '&':
		if cl.tryConsume('&') {
			return tagOpAnd, nil
		}
		if cl.tryConsume('^') {
			return tagOpBitAndNot, nil
		}
		return tagOpBitAnd, nil
	case '|':
		if cl.tryConsume('|') {
			return tagOpOr, nil
		}
		return tagOpBitOr, nil
	case '^':
		return tagOpXor, nil
	case '<':
		if cl.tryConsume('<') {
			return tagOpShl, nil
		}
		if cl.tryConsume('-') {
			return tagOpSend, nil
		}
		if cl.tryConsume('=') {
			return tagOpLte, nil
		}
		return tagOpLt, nil
	case '>':
		if cl.tryConsume('>') {
			return tagOpShr, nil
		}
		if cl.tryConsume('=') {
			return tagOpGte, nil
		}
		return tagOpGt, nil
	case '=':
		if cl.tryConsume('=') {
			return tagOpEq, nil
		}
		return tagOpAssign, nil
	case '!':
		if cl.tryConsume('=') {
			return tagOpNotEq, nil
		}
		return tagOpNot, nil
	case '.':
		return tagMemberAccess, nil

	case scanner.Ident:
		switch s := cl.TokenText(); s {
		case "deref":
			return tagDeref, nil
		case "call":
			return tagCall, nil
		case "callstmt":
			return tagCallStmt, nil
		case "block":
			return tagBlock, nil
		case "array_type":
			return tagArrayType, nil
		case "slice_type":
			return tagSliceType, nil
		}
	}

	return 0, fmt.Errorf("unknown tag: %q", cl.TokenText())
}

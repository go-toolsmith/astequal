package astexp

import (
	"go/ast"
	"go/token"
	"log"
)

// patternTag specifies Pattern category.
type patternTag int

const (
	tagNumber = patternTag(^uint(0)>>1) - iota
	tagCall
	tagCallStmt // Matches ExprStmt{X: CallExpr}
	tagDeref
	tagBlock
	tagMemberAccess
	tagEllipsis

	tagIdent  = patternTag(token.IDENT)
	tagInt    = patternTag(token.INT)
	tagFloat  = patternTag(token.FLOAT)
	tagImag   = patternTag(token.IMAG)
	tagChar   = patternTag(token.CHAR)
	tagString = patternTag(token.STRING)

	tagOpAdd       = patternTag(token.ADD)
	tagOpSub       = patternTag(token.SUB)
	tagOpMul       = patternTag(token.MUL)
	tagOpQuo       = patternTag(token.QUO)
	tagOpRem       = patternTag(token.REM)
	tagOpBitAnd    = patternTag(token.AND)
	tagOpBitOr     = patternTag(token.OR)
	tagOpXor       = patternTag(token.XOR)
	tagOpShl       = patternTag(token.SHL)
	tagOpShr       = patternTag(token.SHR)
	tagOpBitAndNot = patternTag(token.AND_NOT)
	tagOpAnd       = patternTag(token.LAND)
	tagOpOr        = patternTag(token.LOR)
	tagOpSend      = patternTag(token.ARROW)
	tagOpInc       = patternTag(token.INC)
	tagOpDec       = patternTag(token.DEC)
	tagOpEq        = patternTag(token.EQL)
	tagOpLt        = patternTag(token.LSS)
	tagOpGt        = patternTag(token.GTR)
	tagOpAssign    = patternTag(token.ASSIGN)
	tagOpNot       = patternTag(token.NOT)
	tagOpNotEq     = patternTag(token.NEQ)
	tagOpLte       = patternTag(token.LEQ)
	tagOpGte       = patternTag(token.GEQ)
)

// Pattern is compiled Go AST template.
// Can be used for multiple find/match operations.
type Pattern struct {
	tag   patternTag
	value string
	tail  []Pattern
}

// Compile returns Pattern that is described by sexp template.
//
// Internally creates Compiler object for single compilation.
func Compile(sexp string) (*Pattern, error) {
	return NewCompiler().CompileString(sexp)
}

// Match reports whether the Pattern matches the node.
func (pat *Pattern) Match(node ast.Node) bool {
	return pat.match(node)
}

// isFullMatch returns true if given subpat can be interpreted as matched.
func isFullMatch(subpat []Pattern) bool {
	return len(subpat) == 0 ||
		(len(subpat) == 1 && subpat[0].tag == tagEllipsis)
}

// matchStmtList tries to match each ast.Stmt from nodes against subpat.
func (pat *Pattern) matchStmtList(subpat []Pattern, nodes []ast.Stmt) bool {
	for i := range nodes {
		switch {
		case len(subpat) == 0:
			return false
		case subpat[0].tag == tagEllipsis:
			if len(subpat) == 1 {
				return true
			}
			if subpat[1].match(nodes[i]) {
				subpat = subpat[2:]
			}
		default:
			if !subpat[0].match(nodes[i]) {
				return false
			}
			subpat = subpat[1:]
		}
	}
	return isFullMatch(subpat)
}

// matchStmtList tries to match each ast.Expr from nodes against subpat.
func (pat *Pattern) matchExprList(subpat []Pattern, nodes []ast.Expr) bool {
	for i := range nodes {
		switch {
		case len(subpat) == 0:
			return false
		case subpat[0].tag == tagEllipsis:
			if len(subpat) == 1 {
				return true
			}
			if subpat[1].match(nodes[i]) {
				subpat = subpat[2:]
			}
		default:
			if !subpat[0].match(nodes[i]) {
				return false
			}
			subpat = subpat[1:]
		}
	}
	return isFullMatch(subpat)
}

func (pat *Pattern) matchCallExpr(node *ast.CallExpr) bool {
	if !pat.tail[0].match(node.Fun) {
		return false
	}
	return pat.matchExprList(pat.tail[1:], node.Args)
}

func (pat *Pattern) match(node ast.Node) bool {
	switch node := node.(type) {
	case *ast.BlockStmt:
		if pat.tag != tagBlock {
			return false
		}
		return pat.matchStmtList(pat.tail, node.List)

	case *ast.BasicLit:
		if node.Kind == token.Token(pat.tag) && node.Value == pat.value {
			return true
		}

		switch node.Kind {
		case token.INT:
			return pat.tag == tagNumber && node.Value == pat.value
		case token.FLOAT:
			return pat.tag == tagNumber && node.Value == pat.value
		}
		return false

	case *ast.ExprStmt:
		if node, ok := node.X.(*ast.CallExpr); ok {
			switch pat.tag {
			case tagCallStmt, tagCall:
				return pat.matchCallExpr(node)
			default:
				return false
			}
		}
		return pat.match(node.X)

	case *ast.Ident:
		return pat.tag == tagIdent && pat.value == node.Name

	case *ast.CallExpr:
		return pat.tag == tagCall && pat.matchCallExpr(node)

	case *ast.ParenExpr:
		return pat.match(node.X)

	case *ast.BinaryExpr:
		return node.Op == token.Token(pat.tag) &&
			pat.tail[0].match(node.X) &&
			pat.tail[1].match(node.Y)

	case *ast.SelectorExpr:
		return pat.tag == tagMemberAccess &&
			pat.tail[0].match(node.X) &&
			pat.tail[1].match(node.Sel)

	}

	// FIXME: this should be removed when all AST nodes are handled.
	log.Printf("unexpected %T", node)
	return false
}

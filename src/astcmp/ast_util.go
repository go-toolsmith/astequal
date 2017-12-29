package astcmp

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// Variables that are re-used by astParseExpr/astParseStmt/astParseDecl functions.
var (
	badExpr      = &ast.BadExpr{}
	badStmt      = &ast.BadStmt{}
	badDecl      = &ast.BadDecl{}
	emptyFileSet = token.NewFileSet()
)

// astParseExpr parses single expression node from s.
// In case of parse error, ast.BadExpr is returned.
func astParseExpr(s string) ast.Expr {
	node, err := parser.ParseExpr(s)
	if err != nil {
		return badExpr
	}
	return node
}

// astParseStmt parses single statement node from s.
// In case of parse error, ast.BadStmt is returned.
func astParseStmt(s string) ast.Stmt {
	node, err := parser.ParseFile(
		emptyFileSet,
		"",
		"package main;func main() {"+s+"}",
		0)
	if err != nil {
		return badStmt
	}
	return node.Decls[0].(*ast.FuncDecl).Body.List[0]
}

// astParseDecl parses single declaration node from s.
// In case of parse error, ast.BadDecl is returned.
func astParseDecl(s string) ast.Decl {
	node, err := parser.ParseFile(
		emptyFileSet,
		"",
		"package main;"+s,
		0)
	if err != nil {
		return badDecl
	}
	return node.Decls[0]
}

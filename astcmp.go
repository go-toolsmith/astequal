// Package astcmp provides functions to compare AST nodes by their printed
// representations (strings) or directly.
package astcmp

import (
	"go/ast"
)

// Equal returns true for structurally equal AST nodes x and y.
//
// Nil arguments are permitted: true is returned if x and y are both nils.
//
// See also: EqualExpr, EqualStmt, EqualDecl.
func Equal(x, y ast.Node) bool {
	return astNodeEq(x, y)
}

// EqualExpr returns true for structually equal AST expressions x and y.
//
// Nil arguments are permitted: true is returned if x and y are both nils.
// ast.BadExpr comparison always yields false.
func EqualExpr(x, y ast.Expr) bool {
	return astExprEq(x, y)
}

// EqualStmt returns true for structually equal AST statements x and y.
//
// Nil arguments are permitted: true is returned if x and y are both nils.
// ast.BadStmt comparison always yields false.
func EqualStmt(x, y ast.Stmt) bool {
	return astStmtEq(x, y)
}

// EqualDecl returns true for structually equal AST declarations x and y.
//
// Nil arguments are permitted: true is returned if x and y are both nils.
// ast.BadDecl comparison always yields false.
func EqualDecl(x, y ast.Decl) bool {
	return astDeclEq(x, y)
}

// EqualExprString reports if x and y strings represent equivalent expressions.
// Invalid parse result replaced with ast.BadExpr node.
func EqualExprString(x, y string) bool {
	return EqualExpr(astParseExpr(x), astParseExpr(y))
}

// EqualStmtString reports if x and y strings represent equivalent statements.
// Invalid parse result replaced with ast.BadStmt node.
func EqualStmtString(x, y string) bool {
	return EqualStmt(astParseStmt(x), astParseStmt(y))
}

// EqualDeclString reports if x and y strings represent equal declarations.
// Invalid parse result replaced with ast.BadDecl node.
func EqualDeclString(x, y string) bool {
	return EqualDecl(astParseDecl(x), astParseDecl(y))
}

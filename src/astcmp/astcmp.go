package astcmp

import (
	"go/ast"
)

// Equal returns true for structurally equal AST nodes x and y.
// Nil arguments are permitted.
// See also: EqualExpr, EqualStmt, EqualDecl.
func Equal(x, y ast.Node) bool {
	return astNodeEq(x, y)
}

// EqualExpr returns true for structually equal AST expressions x and y.
// Nil arguments are permitted.
// ast.BadExpr comparison always yields false.
func EqualExpr(x, y ast.Expr) bool {
	return astExprEq(x, y)
}

// EqualStmt returns true for structually equal AST statements x and y.
// Nil arguments are permitted.
// ast.BadStmt comparison always yields false.
func EqualStmt(x, y ast.Stmt) bool {
	return astStmtEq(x, y)
}

// EqualDecl returns true for structually equal AST declarations x and y.
// Nil arguments are permitted.
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

// ParseTemplate creates new Pattern that can match specified template.
func ParseTemplate(tmpl string) (*Pattern, error) {
	return nil, nil
}

// Match represents successful pattern match with
// populated named sub-groups (if any).
type Match struct {
	// The matched node itself.
	// Guaranteed to be a parent of any node found in sub-groups.
	Root ast.Node

	// Maps sub-pattern name to the node that matched it.
	// Nil patterns without named captures.
	Captures map[string]ast.Node
}

// Pattern that can be used to match AST nodes.
type Pattern struct {
	matcher matcher
}

// Capture returns node that matched sub-pattern with specified name
// during the last successful match.
//
// If there are more than one pattern with specified name,
// only the one that matched first is accessible.
func (p *Pattern) Capture(name string) ast.Node {
	// What if user wants to disambiguate between "no match" & "nil match"?
	return p.matcher.captures[name]
}

// Match returns true if pattern p matches AST node x.
// Successful match overrides pattern captures.
func (p *Pattern) Match(x ast.Node) bool {
	switch x := x.(type) {
	case ast.Expr:
		return p.matcher.MatchExpr(x)
	case ast.Stmt:
		return p.matcher.MatchStmt(x)
	case ast.Decl:
		return p.matcher.MatchDecl(x)
	default:
		return false
	}
}

// MatchExpr returns true if pattern p matches AST expression x.
// Successful match overrides pattern captures.
func (p *Pattern) MatchExpr(x ast.Expr) bool {
	return p.matcher.MatchExpr(x)
}

// MatchStmt returns true if pattern p matches AST statement x.
// Successful match overrides pattern captures.
func (p *Pattern) MatchStmt(x ast.Stmt) bool {
	return p.matcher.MatchStmt(x)
}

// MatchDecl returns true if pattern p matches AST declaration x.
// Successful match overrides pattern captures.
func (p *Pattern) MatchDecl(x ast.Decl) bool {
	return p.matcher.MatchDecl(x)
}

package astcmp

import (
	"go/ast"
)

// matcher implements AST pattern matching.
type matcher struct {
	// nil for patterns without named sub-patterns.
	captures map[string]ast.Node

	// String values, read-only storage.
	// Shared between all patterns.
	values []string

	// Main pattern and all sub-patterns.
	// Length is at least 1 and first element is main pattern.
	patterns []pattern
}

func (m *matcher) MatchExpr(x ast.Expr) bool {
	return m.matchExpr(&m.patterns[0], x)
}

func (m *matcher) MatchStmt(x ast.Stmt) bool {
	return m.matchStmt(&m.patterns[0], x)
}

func (m *matcher) MatchDecl(x ast.Decl) bool {
	return m.matchDecl(&m.patterns[0], x)
}

func (m *matcher) matchExpr(pat *pattern, x ast.Expr) bool {
	return false // TODO
}

func (m *matcher) matchStmt(pat *pattern, x ast.Stmt) bool {
	return false // TODO
}

func (m *matcher) matchDecl(pat *pattern, x ast.Decl) bool {
	return false // TODO
}

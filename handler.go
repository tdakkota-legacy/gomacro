package macro

import (
	"go/ast"
)

// Handler is macro handler.
type Handler interface {
	Handle(cursor Context, node ast.Node) error
}

// HandlerFunc is type adapter for Handler interface.
type HandlerFunc func(cursor Context, node ast.Node) error

func (f HandlerFunc) Handle(cursor Context, node ast.Node) error {
	return f(cursor, node)
}

// OnlyFunction filter function call nodes using a name.
func OnlyFunction(name string, cb func(ctx Context, call *ast.CallExpr) error) Handler {
	return HandlerFunc(func(cursor Context, node ast.Node) error {
		if callExpr, ok := node.(*ast.CallExpr); ok {
			if f, ok := callExpr.Fun.(*ast.Ident); ok && f.Name == name {
				return cb(cursor, callExpr)
			}
		}

		return nil
	})
}

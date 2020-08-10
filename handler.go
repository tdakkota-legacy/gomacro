package macro

import (
	"go/ast"

	"github.com/tdakkota/gomacro/macroctx"
)

type Handler interface {
	Handle(cursor macroctx.Context, node ast.Node) error
}

type HandlerFunc func(cursor macroctx.Context, node ast.Node) error

func (f HandlerFunc) Handle(cursor macroctx.Context, node ast.Node) error {
	return f(cursor, node)
}

func OnlyFunction(name string, cb func(ctx macroctx.Context, call *ast.CallExpr) error) Handler {
	return HandlerFunc(func(cursor macroctx.Context, node ast.Node) error {
		if callExpr, ok := node.(*ast.CallExpr); ok {
			if f, ok := callExpr.Fun.(*ast.Ident); ok && f.Name == name {
				return cb(cursor, callExpr)
			}
		}

		return nil
	})
}

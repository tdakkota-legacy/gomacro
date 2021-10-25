package deepcopy

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"strings"

	builders "github.com/tdakkota/astbuilders"
)

func target(pkg *types.Package, typ types.Type) *types.Interface {
	params := types.NewTuple()
	result := types.NewVar(0, pkg, "r", typ)
	results := types.NewTuple(result)
	sig := types.NewSignature(nil, params, results, false)
	methods := []*types.Func{
		types.NewFunc(0, pkg, "Copy", sig),
	}
	ityp := types.NewInterfaceType(methods, nil)
	return ityp.Complete()
}

func createFunction(name string, typ ast.Expr, recv *ast.Ident, f builders.BodyFunc) builders.FunctionBuilder {
	return builders.NewFunctionBuilder(name).
		Recv(&ast.Field{
			Names: []*ast.Ident{recv},
			Type:  typ,
		}).
		AddResults([]*ast.Field{
			{
				Names: []*ast.Ident{ast.NewIdent("r")},
				Type:  typ,
			},
		}...).
		Body(f)
}

func changeSelectorHead(sel ast.Expr, n *ast.Ident) (ast.Expr, int, error) { // nolint:gocognit,gocyclo
	indexes := 0
	iter := sel
	stack := []ast.Expr{iter}
loop:
	for {
		switch e := iter.(type) {
		case *ast.SelectorExpr:
			iter = e.X
			stack = append(stack, e.Sel)
		case *ast.IndexExpr:
			iter = e.X
			indexes++
			stack = append(stack, e)
		case *ast.CallExpr:
			iter = e.Fun
			stack = append(stack, e)
		case *ast.ParenExpr:
			iter = e.X
			stack = append(stack, e)
		case *ast.StarExpr:
			iter = e.X
			stack = append(stack, e)
		case *ast.Ident:
			break loop
		}
	}
	var expr ast.Expr = n
	for i := len(stack) - 1; i > 0; i-- {
		switch e := stack[i].(type) {
		case *ast.Ident:
			expr = &ast.SelectorExpr{
				X:   expr,
				Sel: e,
			}
		case *ast.SelectorExpr:
			expr = &ast.SelectorExpr{
				X:   expr,
				Sel: e.Sel,
			}
		case *ast.IndexExpr:
			expr = &ast.IndexExpr{
				X:     expr,
				Index: e.Index,
			}
		case *ast.CallExpr:
			expr = &ast.CallExpr{
				Fun:  expr,
				Args: e.Args,
			}
		case *ast.ParenExpr:
			expr = &ast.ParenExpr{
				X: expr,
			}
		case *ast.StarExpr:
			expr = &ast.StarExpr{
				X: expr,
			}
		default:
			var s strings.Builder
			_ = printer.Fprint(&s, token.NewFileSet(), stack[i])
			return nil, 0, fmt.Errorf("unexpected type %T: %s", stack[i], s.String())
		}
	}

	return expr, indexes, nil
}

func elemType(pkg *types.Package, elem types.Type) ast.Expr {
	typ := types.TypeString(elem, func(i *types.Package) string {
		if i.Path() != pkg.Path() {
			return i.Name()
		}
		return ""
	})
	split := strings.Split(typ, ".")
	if len(split) > 1 {
		return builders.SelectorName(split[0], split[1], split[2:]...)
	}
	return ast.NewIdent(split[0])
}

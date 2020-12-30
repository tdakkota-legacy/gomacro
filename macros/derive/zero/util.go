package main

import (
	"go/ast"
	"go/types"
	"strings"

	builders "github.com/tdakkota/astbuilders"
)

func target(pkg *types.Package) *types.Interface {
	params := types.NewTuple()
	result := types.NewVar(0, pkg, "r", types.Typ[types.Bool])
	results := types.NewTuple(result)
	sig := types.NewSignature(nil, params, results, false)
	methods := []*types.Func{
		types.NewFunc(0, pkg, "Zero", sig),
	}
	ityp := types.NewInterfaceType(methods, nil)
	return ityp.Complete()
}

func createFunction(name string, typ ast.Expr, bodyFunc func(*ast.Ident, builders.StatementBuilder) builders.StatementBuilder) builders.FunctionBuilder {
	recv := ast.NewIdent("m")
	return builders.NewFunctionBuilder(name).
		Recv(&ast.Field{
			Names: []*ast.Ident{recv},
			Type:  typ,
		}).
		AddResults([]*ast.Field{
			{
				Names: []*ast.Ident{ast.NewIdent("r")},
				Type:  ast.NewIdent("bool"),
			},
		}...).
		Body(func(s builders.StatementBuilder) builders.StatementBuilder {
			return bodyFunc(recv, s)
		})
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

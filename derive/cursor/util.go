package cursor

import (
	"errors"
	"go/ast"
	"go/types"
	"os"

	builders "github.com/tdakkota/astbuilders"

	"golang.org/x/tools/go/packages"
)

const pkg = "github.com/tdakkota/cursor"

func CreateFunction(name string, typ ast.Expr, bodyFunc builders.BodyFunc) builders.FunctionBuilder {
	selector := ast.NewIdent("m")
	return builders.NewFunctionBuilder(name).
		Recv(&ast.Field{
			Names: []*ast.Ident{selector},
			Type:  typ,
		}).
		AddParameters([]*ast.Field{
			{
				Names: []*ast.Ident{ast.NewIdent("cur")},
				Type:  builders.RefFor(builders.SelectorName("cursor", "Cursor")),
			},
		}...).
		AddResults([]*ast.Field{
			{
				Names: []*ast.Ident{ast.NewIdent("err")},
				Type:  ast.NewIdent("error"),
			},
		}...).
		Body(bodyFunc)
}

var ErrFailedToFindCursor = errors.New("failed to import cursor package")

func load(pkg string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedImports,
		Env:  os.Environ(),
	}
	return packages.Load(cfg, pkg)
}

func target(pkgs []*packages.Package, name string) (*types.Interface, error) {
	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(name)
		if obj == nil {
			continue
		}

		i, ok := obj.Type().(*types.Named)
		if !ok {
			return nil, ErrFailedToFindCursor
		}

		return i.Underlying().(*types.Interface), nil
	}

	return nil, ErrFailedToFindCursor
}

func checkErr(s builders.StatementBuilder) builders.StatementBuilder {
	nilIdent := ast.NewIdent("nil")
	errIdent := ast.NewIdent("err")

	cond := builders.NotEq(errIdent, nilIdent)
	return s.If(nil, cond, func(ifBody builders.StatementBuilder) builders.StatementBuilder {
		return ifBody.Return(errIdent)
	})
}

func callCurFunc(selector ast.Expr, name string) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()
	sel := builders.Selector(selector, ast.NewIdent(name))

	s = s.Define(ast.NewIdent("err"))(builders.Call(sel, ast.NewIdent("cur")))
	s = checkErr(s)

	return s.CompleteAsBlock(), nil
}

func elemType(pkg *types.Package, elem types.Type) *ast.Ident {
	typ := types.TypeString(elem, func(i *types.Package) string {
		if i.Path() != pkg.Path() {
			return i.Name()
		}
		return ""
	})
	return ast.NewIdent(typ)
}

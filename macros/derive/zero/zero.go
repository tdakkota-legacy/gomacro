package zero

import (
	"fmt"
	"go/ast"
	"go/types"

	builders "github.com/tdakkota/astbuilders"
	macro "github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/derive"
)

type zeroProtocol struct{}

func (z zeroProtocol) onNonZero() builders.BodyFunc {
	return func(body builders.StatementBuilder) builders.StatementBuilder {
		return body.Return()
	}
}

func (z zeroProtocol) Array(d *derive.Derive, field derive.Field, arr derive.Array) (*ast.BlockStmt, error) {
	if arr.Size >= 0 {
		return z.notNilableType(field, builders.ArrayOfSize(elemType(d.Package, arr.Elem), int(arr.Size)))
	}
	return z.nilableType(field)
}

func (z zeroProtocol) Map(_ *derive.Derive, field derive.Field, _ derive.Map) (*ast.BlockStmt, error) {
	return z.nilableType(field)
}

func (z zeroProtocol) notNilableType(field derive.Field, typ ast.Expr) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()
	zero := ast.NewIdent("zero")
	s = s.Var(&ast.ValueSpec{
		Names: []*ast.Ident{zero},
		Type:  typ,
	})
	s = s.If(nil, builders.NotEq(field.Selector, zero), z.onNonZero())
	return s.CompleteAsBlock(), nil
}

func (z zeroProtocol) nilableType(field derive.Field) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()
	s = s.If(nil, builders.NotEq(field.Selector, builders.Nil()), z.onNonZero())
	return s.CompleteAsBlock(), nil
}

func (z zeroProtocol) CallFor(d *derive.Derive, field derive.Field, kind types.BasicKind) (*ast.BlockStmt, error) {
	if field.Named {
		return z.notNilableType(field, elemType(d.Package, field.TypeName.Type()))
	}
	return z.notNilableType(field, ast.NewIdent(types.Typ[kind].Name()))
}

func (z zeroProtocol) Impl(d *derive.Derive, field derive.Field) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()
	sel := builders.Selector(field.Selector, ast.NewIdent("Zero"))
	s = s.If(nil, builders.Not(builders.Call(sel)), z.onNonZero())
	return s.CompleteAsBlock(), nil
}

func (z zeroProtocol) Callback(d *derive.Derive, typeSpec *ast.TypeSpec) error {
	if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
		return nil
	}

	obj := d.TypesInfo.ObjectOf(typeSpec.Name)
	if obj == nil {
		return fmt.Errorf("failed to load type info for %s", typeSpec.Name.Name)
	}

	var err error
	recv := ast.NewIdent("m")
	builder := createFunction("Zero", recv, typeSpec.Name,
		func(s builders.StatementBuilder) builders.StatementBuilder {
			if types.Comparable(obj.Type()) {
				lit := &ast.CompositeLit{
					Type: typeSpec.Name,
				}
				return s.Return(builders.Eq(recv, lit))
			}

			s, err = d.Derive(typeSpec, s)
			return s.Return(ast.NewIdent("true"))
		})

	d.AddDecls(builder.CompleteAsDecl())
	return err
}

func Zero(name string) macro.Handler {
	return macro.HandlerFunc(func(cursor macro.Context, node ast.Node) error {
		if cursor.Pre { // skip first pass
			return nil
		}

		proto := &zeroProtocol{}
		zero := derive.CreateMacro(name, target(cursor.Package), proto)
		return derive.Callback(zero).Handle(cursor, node)
	})
}

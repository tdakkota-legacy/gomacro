package cursor

import (
	"go/ast"
	"go/types"
	"strings"

	builders "github.com/tdakkota/astbuilders"
	"github.com/tdakkota/gomacro/derive"
)

type Serialize struct{}

func (m *Serialize) CallFor(d *derive.Derive, field derive.Field, kind types.BasicKind) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()
	name := "Write" + strings.Title(types.Typ[kind].String())

	results := []ast.Expr{builders.Err()}
	if kind == types.String {
		results = append([]ast.Expr{ast.NewIdent("_")}, results...)
	}

	sel := field.Selector
	if field.TypeName != nil {
		sel = builders.Cast(builders.IdentOfKind(kind), sel)
	}

	s = s.Define(results[0], results[1:]...)(builders.CastPackage("cur", name, sel))
	s = checkErr(s)

	return s.CompleteAsBlock(), nil
}

func (m *Serialize) Array(d *derive.Derive, field derive.Field, arr derive.Array) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()

	if v, ok := field.Tag.Lookup("length"); ok {
		_, err := d.Interpolator.ExprExpectInfo(v, types.IsInteger)
		if err != nil {
			return nil, err
		}
	} else if arr.Size <= -1 {
		length := builders.Cast(builders.IdentOfKind(types.Uint8), builders.Len(field.Selector))
		s = s.Define(builders.Err())(builders.CallPackage("cur", "WriteUint8", length))
		s = checkErr(s)
	}

	var err error
	value := ast.NewIdent("v")
	s = s.Range(ast.NewIdent("_"), value, field.Selector,
		func(loop builders.StatementBuilder) builders.StatementBuilder {
			loop, err = d.Dispatch(derive.Field{
				Selector: value,
			}, arr.Elem, loop)
			return loop
		})

	return s.CompleteAsBlock(), err
}

func (m *Serialize) Impl(d *derive.Derive, field derive.Field) (*ast.BlockStmt, error) {
	return callCurFunc(field.Selector, "Append")
}

func (m *Serialize) Callback(d *derive.Derive, typeSpec *ast.TypeSpec) error {
	if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
		return nil
	}

	var err error
	builder := CreateFunction("Append", typeSpec.Name, func(s builders.StatementBuilder) builders.StatementBuilder {
		s, err = d.Derive(typeSpec, s)
		return s.Return(builders.Nil())
	})

	d.AddDecls(builder.CompleteAsDecl())
	return err
}

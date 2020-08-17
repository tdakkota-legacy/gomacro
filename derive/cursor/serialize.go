package cursor

import (
	"go/ast"
	"go/types"
	"strings"

	builders "github.com/tdakkota/astbuilders"
	"github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/derive"
	"github.com/tdakkota/gomacro/derive/base"
)

type Serialize struct {
	Derive *derive.Derive
}

func NewSerialize() *Serialize {
	return &Serialize{Derive: nil}
}

func (m *Serialize) CallFor(field base.Field, kind types.BasicKind) (*ast.BlockStmt, error) {
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

func (m *Serialize) Array(d base.Dispatcher, field base.Field, arr derive.Array) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()

	if v, ok := field.Tag.Lookup("length"); ok {
		_, err := m.Derive.Interpolator.ExprExpectInfo(v, types.IsInteger)
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
			loop, err = d.Dispatch(base.Field{
				Selector: value,
			}, arr.Elem, loop)
			return loop
		})

	return s.CompleteAsBlock(), err
}

func (m *Serialize) Impl(field base.Field) (*ast.BlockStmt, error) {
	return callCurFunc(field.Selector, "Append")
}

func (m *Serialize) create(context macro.Context) error {
	if m.Derive != nil {
		return nil
	}

	i, err := target("Appender")
	if err != nil {
		return err
	}

	info := derive.NewDeriveInfo(m, "derive_binary", i)
	m.Derive = derive.NewDerive(context, info)

	return nil
}

func (m *Serialize) Callback(context macro.Context, node ast.Node) error {
	if typeSpec, ok := node.(*ast.TypeSpec); ok {
		if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
			return nil
		}

		err := m.create(context)
		if err != nil {
			return err
		}

		builder := CreateFunction("Append", typeSpec.Name, func(s builders.StatementBuilder) builders.StatementBuilder {
			s, err = m.Derive.Derive(typeSpec, s)
			return s.Return(builders.Nil())
		})

		context.AddDecls(builder.CompleteAsDecl())
		return err
	}

	return nil
}

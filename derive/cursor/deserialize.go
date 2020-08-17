package cursor

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	builders "github.com/tdakkota/astbuilders"
	"github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/derive"
	"github.com/tdakkota/gomacro/derive/base"
)

type Deserialize struct {
	Derive *derive.Derive
}

func NewDeserialize() *Deserialize {
	return &Deserialize{Derive: nil}
}

func (m *Deserialize) CallFor(field base.Field, kind types.BasicKind) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()
	name := "Read" + strings.Title(types.Typ[kind].String())

	var tmp ast.Expr = ast.NewIdent("tmp")
	s = s.Define(tmp, builders.Err())(builders.CallPackage("cur", name))
	s = checkErr(s)

	if field.TypeName != nil {
		if m.Derive.Package.Path() != field.TypeName.Pkg().Path() {
			tmp = builders.CastPackage(field.TypeName.Pkg().Name(), field.TypeName.Name(), tmp)
		} else {
			tmp = builders.Cast(ast.NewIdent(field.TypeName.Name()), tmp)
		}
	}

	sel := field.Selector
	if _, ok := sel.(*ast.Ident); ok {
		sel = builders.DeRef(sel)
	}
	s = s.Assign(sel)(token.ASSIGN)(tmp)

	return s.CompleteAsBlock(), nil
}

func (m *Deserialize) Array(d base.Dispatcher, field base.Field, arr derive.Array) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()
	size := ast.NewIdent("n")

	if v, ok := field.Tag.Lookup("length"); ok {
		expr, err := m.Derive.Interpolator.ExprExpectInfo(v, types.IsInteger)
		if err != nil {
			return nil, err
		}

		s = s.Define(size)(expr)
	} else {
		if arr.Size <= -1 {
			s = s.Define(size, builders.Err())(builders.CallPackage("cur", "ReadUint8"))
			s = checkErr(s)
		} else {
			s = s.Define(size)(builders.IntegerLit(int(arr.Size)))
		}
	}

	i := ast.NewIdent("i")
	init := builders.Assign(i)(token.DEFINE)(builders.IntegerLit(0))
	inc := &ast.IncDecStmt{X: i, Tok: token.INC}

	var err error
	to := builders.Cast(ast.NewIdent("int"), size)
	s = s.For(init, builders.Less(i, to), inc, func(loop builders.StatementBuilder) builders.StatementBuilder {
		loop, err = d.Dispatch(base.Field{
			Selector: &ast.IndexExpr{
				X:     field.Selector,
				Index: i,
			},
		}, arr.Elem, loop)
		return loop
	})

	return s.CompleteAsBlock(), err
}

func (m *Deserialize) Impl(field base.Field) (*ast.BlockStmt, error) {
	return callCurFunc(field.Selector, "Scan")
}

func (m *Deserialize) create(context macro.Context) error {
	if m.Derive != nil {
		return nil
	}

	i, err := target("Scanner")
	if err != nil {
		return err
	}

	info := derive.NewDeriveInfo(m, "derive_binary", i)
	m.Derive = derive.NewDerive(context, info)

	return nil
}

func (m *Deserialize) Callback(context macro.Context, node ast.Node) error {
	if typeSpec, ok := node.(*ast.TypeSpec); ok {
		if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
			return nil
		}

		err := m.create(context)
		if err != nil {
			return err
		}

		builder := CreateFunction("Scan", builders.RefFor(typeSpec.Name), func(s builders.StatementBuilder) builders.StatementBuilder {
			s, err = m.Derive.Derive(typeSpec, s)
			return s.Return(builders.Nil())
		})

		context.AddDecls(builder.CompleteAsDecl())
		return err
	}

	return nil
}

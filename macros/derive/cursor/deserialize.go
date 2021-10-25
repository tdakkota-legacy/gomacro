package cursor

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	builders "github.com/tdakkota/astbuilders"
	"github.com/tdakkota/gomacro/derive"
)

// Deserialize defines derive.Protocol for binary deserialization.
type Deserialize struct{}

// CallFor implements derive.Protocol.
func (m *Deserialize) CallFor(d *derive.Derive, field derive.Field, kind types.BasicKind) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()
	name := "Read" + strings.Title(types.Typ[kind].String())

	var tmp ast.Expr = ast.NewIdent("tmp")
	s = s.Define(tmp, builders.Err())(builders.CallPackage("cur", name))
	s = checkErr(s)

	if field.TypeName != nil {
		if d.Package.Path() != field.TypeName.Pkg().Path() {
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

func (m *Deserialize) createArray(size, sel, typ ast.Expr, s builders.StatementBuilder) builders.StatementBuilder {
	return s.Assign(sel)(token.ASSIGN)(builders.MakeExpr(builders.SliceOf(typ), size, nil))
}

// Array implements derive.ArrayDerive.
func (m *Deserialize) Array(d *derive.Derive, field derive.Field, arr derive.Array) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()
	size := ast.NewIdent("n")

	elemTyp := elemType(d.Package, arr.Elem)
	if v, ok := field.Tag.Lookup("length"); ok {
		expr, err := d.Interpolator.ExprExpectInfo(v, types.IsInteger)
		if err != nil {
			return nil, err
		}

		s = s.Define(size)(expr)
		s = m.createArray(size, field.Selector, elemTyp, s)
	} else {
		if arr.Size <= -1 {
			s = s.Define(size, builders.Err())(builders.CallPackage("cur", "ReadUint8"))
			s = checkErr(s)
			s = m.createArray(size, field.Selector, elemTyp, s)
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
		loop, err = d.Dispatch(derive.Field{
			Selector: &ast.IndexExpr{
				X:     field.Selector,
				Index: i,
			},
		}, arr.Elem, loop)
		return loop
	})

	return s.CompleteAsBlock(), err
}

// Impl implements derive.Protocol.
func (m *Deserialize) Impl(d *derive.Derive, field derive.Field) (*ast.BlockStmt, error) {
	return callCurFunc(field.Selector, "Scan")
}

// Callback implements derive.Protocol.
func (m *Deserialize) Callback(d *derive.Derive, typeSpec *ast.TypeSpec) error {
	if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
		return nil
	}

	var err error
	builder := createFunction("Scan", builders.RefFor(typeSpec.Name), func(s builders.StatementBuilder) builders.StatementBuilder {
		s, err = d.Derive(typeSpec, s)
		return s.Return(builders.Nil())
	})

	d.AddDecls(builder.CompleteAsDecl())
	return err
}

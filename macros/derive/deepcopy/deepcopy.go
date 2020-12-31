package deepcopy

import (
	"fmt"
	builders "github.com/tdakkota/astbuilders"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	macro "github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/derive"
)

type deepCopyProtocol struct {
}

func (p deepCopyProtocol) tmpAndRef(to, from ast.Expr, s builders.StatementBuilder) builders.StatementBuilder {
	tmp := ast.NewIdent("tmp")
	s = s.Define(tmp)(from)
	return s.Assign(to)(token.ASSIGN)(builders.Ref(tmp))
}

func (p deepCopyProtocol) comparable(field derive.Field) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()
	expr, _, err := changeSelectorHead(field.Selector, ast.NewIdent("r"))
	if err != nil {
		return nil, err
	}
	s = s.Assign(expr)(token.ASSIGN)(field.Selector)
	return s.CompleteAsBlock(), nil
}

func (p deepCopyProtocol) Map(d *derive.Derive, field derive.Field, m derive.Map) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()
	expr, idx, err := changeSelectorHead(field.Selector, ast.NewIdent("r"))
	if err != nil {
		return nil, err
	}

	makeExpr := builders.MakeExpr(
		builders.MapOf(elemType(d.Package, m.Key), elemType(d.Package, m.Value)),
		builders.Len(field.Selector), nil,
	)
	s = s.Assign(expr)(token.ASSIGN)(makeExpr)

	index := ast.NewIdent(strings.Repeat("k", idx+1))
	if st, ok := m.Value.(*types.Struct); ok && st.NumFields() == 0 {
		s = s.Range(index, nil, field.Selector, func(loop builders.StatementBuilder) builders.StatementBuilder {
			indx := builders.Index(expr, index)
			return loop.Assign(indx)(token.ASSIGN)(&ast.CompositeLit{
				Type: builders.EmptyStruct(),
			})
		})
	} else {
		s = s.Range(index, nil, field.Selector, func(loop builders.StatementBuilder) builders.StatementBuilder {
			loop, err = d.Dispatch(derive.Field{
				Selector: builders.Index(field.Selector, index),
			}, m.Value, loop)
			return loop
		})
	}

	if err != nil {
		return nil, err
	}

	return s.CompleteAsBlock(), nil
}

func (p deepCopyProtocol) Pointer(d *derive.Derive, field derive.Field, ptr derive.Pointer) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()
	expr, _, err := changeSelectorHead(field.Selector, ast.NewIdent("r"))
	if err != nil {
		return nil, err
	}

	var body builders.BodyFunc
	if named, ok := ptr.Elem.(*types.Named); ok && d.IsDelayed(named.Obj()) {
		// Call Copy func if this type have deepcopy derive annotation too.
		body = func(ifBody builders.StatementBuilder) builders.StatementBuilder {
			sel := builders.Selector(field.Selector, ast.NewIdent("Copy"))
			return p.tmpAndRef(expr, builders.Call(sel), ifBody)
		}
	} else if !types.Comparable(ptr.Elem) {
		// Copy field-by-field if it contains pointers
		s = p.tmpAndRef(expr, builders.DeRef(field.Selector), s)

		body = func(ifBody builders.StatementBuilder) builders.StatementBuilder {
			ifBody, err = d.Dispatch(derive.Field{
				Selector: field.Selector,
			}, ptr.Elem, ifBody)
			return ifBody
		}
	} else {
		body = func(ifBody builders.StatementBuilder) builders.StatementBuilder {
			return ifBody.Assign(expr)(token.ASSIGN)(builders.DeRef(field.Selector))
		}
	}

	s = s.If(nil, builders.NotEq(field.Selector, builders.Nil()), body)
	if err != nil {
		return nil, err
	}

	return s.CompleteAsBlock(), nil
}

func (p deepCopyProtocol) Array(d *derive.Derive, field derive.Field, arr derive.Array) (*ast.BlockStmt, error) {
	if arr.Size >= 0 { // array case
		return p.comparable(field)
	}

	s := builders.NewStatementBuilder()
	expr, idx, err := changeSelectorHead(field.Selector, ast.NewIdent("r"))
	if err != nil {
		return nil, err
	}

	makeExpr := builders.MakeExpr(
		builders.SliceOf(elemType(d.Package, arr.Elem)),
		builders.Len(field.Selector), nil,
	)
	s = s.Assign(expr)(token.ASSIGN)(makeExpr)

	// Here we check the slice element is comparable.
	// Maps, slices, chans or structs which contains at least one of them
	// are not comparable, so are not copyable via copy builtin,
	// because this types should be copied explicitly.
	if !types.Comparable(arr.Elem) {
		index := ast.NewIdent(strings.Repeat("i", idx+1))
		s = s.Range(index, nil, field.Selector,
			func(loop builders.StatementBuilder) builders.StatementBuilder {
				loop, err = d.Dispatch(derive.Field{
					Selector: builders.Index(field.Selector, index),
				}, arr.Elem, loop)
				return loop
			})

		if err != nil {
			return nil, err
		}
	} else {
		s = s.Expr(builders.Copy(expr, field.Selector))
	}

	return s.CompleteAsBlock(), nil
}

func (p deepCopyProtocol) CallFor(d *derive.Derive, field derive.Field, kind types.BasicKind) (*ast.BlockStmt, error) {
	return p.comparable(field)
}

func (p deepCopyProtocol) Impl(d *derive.Derive, field derive.Field) (*ast.BlockStmt, error) {
	s := builders.NewStatementBuilder()
	sel := builders.Selector(field.Selector, ast.NewIdent("Copy"))
	expr, _, err := changeSelectorHead(field.Selector, ast.NewIdent("r"))
	if err != nil {
		return nil, err
	}
	s = s.Assign(expr)(token.ASSIGN)(builders.Call(sel))
	return s.CompleteAsBlock(), nil
}

func (p deepCopyProtocol) Callback(d *derive.Derive, typeSpec *ast.TypeSpec) error {
	if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
		return nil
	}

	obj := d.TypesInfo.ObjectOf(typeSpec.Name)
	if obj == nil {
		return fmt.Errorf("failed to load type info for %s", typeSpec.Name.Name)
	}

	var err error
	recv := ast.NewIdent("m")
	builder := createFunction("Copy", typeSpec.Name, recv,
		func(s builders.StatementBuilder) builders.StatementBuilder {
			s, err = d.Derive(typeSpec, s)
			return s.Return()
		})

	d.AddDecls(builder.CompleteAsDecl())
	return err
}

func DeepCopy(name string) macro.Handler {
	return macro.HandlerFunc(func(cursor macro.Context, node ast.Node) error {
		if cursor.Pre { // skip first pass
			return nil
		}

		typeSpec, ok := node.(*ast.TypeSpec)
		if !ok {
			return nil
		}

		obj := cursor.TypesInfo.ObjectOf(typeSpec.Name)
		if obj == nil {
			return fmt.Errorf("failed to load type info for %s", typeSpec.Name.Name)
		}

		proto := &deepCopyProtocol{}
		zero := derive.CreateMacro(name, target(cursor.Package, obj.Type()), proto)
		return derive.Callback(zero).Handle(cursor, node)
	})
}

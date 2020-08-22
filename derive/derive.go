package derive

import (
	"fmt"
	"go/ast"
	"go/types"

	builders "github.com/tdakkota/astbuilders"
	"github.com/tdakkota/gomacro"
)

type Derive struct {
	macro.Context
	Macro
	Interpolator Interpolator

	delayed  macro.DelayedTypes
	first    bool
	obj      *types.TypeName
	typeSpec *ast.TypeSpec
	selector *ast.Ident
}

func NewDerive(context macro.Context, m Macro) *Derive {
	selector := ast.NewIdent("m")
	d := &Derive{
		Context:  context,
		Macro:    m,
		delayed:  context.Delayed[m.Name()],
		first:    true,
		selector: selector,
	}

	d.Interpolator = NewInterpolator(d, "$m", selector.Name)
	return d
}

//nolint: unparam
func (d *Derive) impl(field Field, typ types.Type, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	call, err := d.Protocol().Impl(d, field)
	if err != nil {
		return s, err
	}

	return s.AddStmts(call), nil
}

func (d *Derive) IsCurrent(typ types.Type) bool {
	return types.AssignableTo(typ, d.TypesInfo.TypeOf(d.typeSpec.Name))
}

func (d *Derive) dispatch1(field Field, typ types.Type, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	// skip tag
	if _, ok := field.Tag.Lookup("skip"); ok {
		return s, nil
	}

	// if tag
	if v, ok := field.Tag.Lookup("if"); ok {
		expr, err := d.Interpolator.ExprExpectInfo(v, types.IsBoolean)
		if err != nil {
			return s, fmt.Errorf("failed to parse expression: %w", err)
		}

		s = s.If(nil, expr, func(body builders.StatementBuilder) builders.StatementBuilder {
			body, err = d.dispatch(field, false, typ, body)
			return body
		})

		return s, err
	}

	return d.dispatch(field, false, typ, s)
}

//nolint:gocyclo
func (d *Derive) dispatch(field Field, named bool, typ types.Type, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	// Types, which will be implemented later
	if field.TypeName != nil && !d.IsCurrent(typ) {
		if v, ok := typ.(container); ok {
			elem := v.Elem()
			if !d.IsCurrent(elem) && checkType(elem, field.TypeName) {
				return s, fmt.Errorf("%w: '%s' refers to itself", ErrCycleDetected, field.TypeName.Name())
			}
		}
	}

	// User-defined target implementations
	if types.Implements(typ, d.Macro.Target()) {
		return d.impl(field, d.Macro.Target(), s)
	}

	switch v := typ.(type) {
	case *types.Basic:
		if !named {
			field.TypeName = nil
		}
		return d.Basic(field, v, s)
	case *types.Array:
		return d.Array(field, v, s)
	case *types.Slice:
		return d.Slice(field, v, s)
	case *types.Struct:
		return d.Struct(field, v, s)
	case *types.Pointer:
		return d.Pointer(field, v, s)
	case *types.Interface:
		return d.Interface(field, v, s)
	case *types.Map:
		if field.TypeName != nil && checkType(v.Key(), field.TypeName) {
			return s, fmt.Errorf("%w: '%s' refers to itself", ErrCycleDetected, field.TypeName.Name())
		}

		return d.Map(field, v, s)
	case *types.Chan:
		return d.Chan(field, v, s)
	case *types.Named:
		field.TypeName = v.Obj()
		if d.delayed.Find(field.TypeName) && !d.first {
			return d.impl(field, d.Macro.Target(), s)
		}
		if d.first {
			d.first = false
		}

		return d.dispatch(field, true, v.Underlying(), s)
	}

	return s, nil
}

func (d *Derive) Dispatch(field Field, typ types.Type, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	return d.dispatch1(field, typ, s)
}

func (d *Derive) Derive(t *ast.TypeSpec, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	d.typeSpec = t

	obj := d.TypesInfo.ObjectOf(d.typeSpec.Name)
	if name, ok := obj.(*types.TypeName); ok {
		d.obj = name
	} else {
		return s, fmt.Errorf("failed to load type info for %s", d.typeSpec.Name.Name)
	}

	field := Field{
		TypeName: d.obj,
		Selector: d.selector,
	}

	return d.dispatch1(field, d.TypesInfo.TypeOf(d.typeSpec.Name), s)
}

func (d *Derive) Basic(field Field, typ *types.Basic, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	if typ.Kind() == types.Invalid {
		return s, fmt.Errorf("%v: %w", field.Selector, ErrInvalidType)
	}

	block, err := d.Protocol().CallFor(d, field, typ.Kind())
	if err != nil {
		return s, err
	}

	return s.AddStmts(block), nil
}

func (d *Derive) array(array Array, field Field, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	if v, ok := d.Macro.Protocol().(ArrayDerive); ok {
		stmts, err := v.Array(d, field, array)
		if err != nil {
			return s, err
		}
		return s.AddStmts(stmts), nil
	}

	var err error
	value := ast.NewIdent("v")
	s = s.Range(ast.NewIdent("_"), value, field.Selector,
		func(loop builders.StatementBuilder) builders.StatementBuilder {
			loop, err = d.dispatch(Field{
				Selector: value,
			}, false, array.Elem, loop)
			return loop
		})
	return s, err
}

func (d *Derive) Array(field Field, typ *types.Array, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	return d.array(Array{typ.Len(), typ.Elem()}, field, s)
}

func (d *Derive) Slice(field Field, typ *types.Slice, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	return d.array(Array{-1, typ.Elem()}, field, s)
}

func (d *Derive) Struct(field Field, typ *types.Struct, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	var err error
	for i := 0; i < typ.NumFields(); i++ {
		subField := typ.Field(i)
		if !subField.Exported() && subField.Pkg().Path() != d.Context.Package.Path() {
			continue
		}

		parentSelector := field.Selector
		if parentSelector == nil {
			parentSelector = d.selector
		}

		newField := Field{
			TypeName: field.TypeName,
			Selector: builders.Selector(parentSelector, ast.NewIdent(subField.Name())),
			Tag:      Tag(typ.Tag(i)),
		}

		s, err = d.dispatch1(newField, subField.Type(), s)
		if err != nil {
			return s, err
		}
	}

	return s, nil
}

// Ignore interface marshaling (types which implements marshalling interface already handled)
func (d *Derive) Interface(field Field, typ *types.Interface, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	return s, nil
}

// Ignore pointer marshaling
func (d *Derive) Pointer(field Field, typ *types.Pointer, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	return s, nil
}

func (d *Derive) Map(field Field, typ *types.Map, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	if v, ok := d.Macro.Protocol().(MapDerive); ok {
		stmts, err := v.Map(d, field, Map{
			Key:   typ.Key(),
			Value: typ.Elem(),
		})
		if err != nil {
			return s, err
		}
		return s.AddStmts(stmts), nil
	}

	return s, nil
}

// Ignore chan marshaling
func (d *Derive) Chan(field Field, typ *types.Chan, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	return s, nil
}

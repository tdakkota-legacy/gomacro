package derive

import (
	"fmt"
	"go/ast"
	"go/types"
	"reflect"

	builders "github.com/tdakkota/astbuilders"
	"github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/derive/base"
)

type Derive struct {
	macro.Context
	Info

	delayed      macro.DelayedTypes
	first        bool
	typeSpec     *ast.TypeSpec
	interpolator interpolator
	selector     *ast.Ident
	arrayBitSize int
}

func NewDerive(context macro.Context, deriveInfo Info) *Derive {
	return &Derive{
		Context:      context,
		Info:         deriveInfo,
		delayed:      context.Delayed[deriveInfo.macroName],
		first:        true,
		interpolator: newInterpolator("$m", "m"),
		selector:     ast.NewIdent("m"),
		arrayBitSize: 8,
	}
}

//nolint: unparam
func (d *Derive) impl(field base.Field, typ types.Type, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	call, err := d.Impl(field)
	if err != nil {
		return s, err
	}

	return s.AddStmts(call), nil
}

func (d *Derive) IsCurrent(typ types.Type) bool {
	return types.AssignableTo(typ, d.TypesInfo.TypeOf(d.typeSpec.Name))
}

func (d *Derive) dispatch1(field base.Field, typ types.Type, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	if v, ok := field.Tag.Lookup("if"); ok {
		expr, err := d.interpolator.Expr(v)
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
func (d *Derive) dispatch(field base.Field, named bool, typ types.Type, s builders.StatementBuilder) (builders.StatementBuilder, error) {
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
	if types.Implements(typ, d.target) {
		return d.impl(field, d.target, s)
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
			return d.impl(field, d.target, s)
		}
		if d.first {
			d.first = false
		}

		return d.dispatch(field, true, v.Underlying(), s)
	}

	return s, nil
}

func (d *Derive) Dispatch(field base.Field, typ types.Type, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	return d.dispatch1(field, typ, s)
}

func (d *Derive) Derive(t *ast.TypeSpec, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	d.typeSpec = t

	field := base.Field{
		TypeName: d.TypesInfo.ObjectOf(d.typeSpec.Name).(*types.TypeName),
		Selector: d.selector,
	}
	return d.dispatch1(field, d.TypesInfo.TypeOf(d.typeSpec.Name), s)
}

func (d *Derive) Basic(field base.Field, typ *types.Basic, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	if typ.Kind() == types.Invalid {
		return s, fmt.Errorf("%v: %w", field.Selector, ErrInvalidType)
	}

	block, err := d.CallFor(field, typ.Kind())
	if err != nil {
		return s, err
	}

	return s.AddStmts(block), nil
}

func (d *Derive) array(array Array, field base.Field, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	if v, ok := d.Info.Interface.(ArrayDerive); ok {
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
			loop, err = d.dispatch(base.Field{
				Selector: value,
			}, false, array.Elem, loop)
			return loop
		})
	return s, err
}

func (d *Derive) Array(field base.Field, typ *types.Array, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	return d.array(Array{typ.Len(), typ.Elem()}, field, s)
}

func (d *Derive) Slice(field base.Field, typ *types.Slice, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	return d.array(Array{-1, typ.Elem()}, field, s)
}

func (d *Derive) Struct(field base.Field, typ *types.Struct, s builders.StatementBuilder) (builders.StatementBuilder, error) {
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

		newField := base.Field{
			TypeName: field.TypeName,
			Selector: builders.Selector(parentSelector, ast.NewIdent(subField.Name())),
			Tag:      reflect.StructTag(typ.Tag(i)),
		}

		s, err = d.dispatch1(newField, subField.Type(), s)
		if err != nil {
			return s, err
		}
	}

	return s, nil
}

// Ignore interface marshaling (types which implements marshalling interface already handled)
func (d *Derive) Interface(field base.Field, typ *types.Interface, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	return s, nil
}

// Ignore pointer marshaling
func (d *Derive) Pointer(field base.Field, typ *types.Pointer, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	return s, nil
}

func (d *Derive) Map(field base.Field, typ *types.Map, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	if v, ok := d.Info.Interface.(MapDerive); ok {
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
func (d *Derive) Chan(field base.Field, typ *types.Chan, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	return s, nil
}

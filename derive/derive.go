package derive

import (
	"fmt"
	"go/ast"
	"go/types"

	builders "github.com/tdakkota/astbuilders"
	macro "github.com/tdakkota/gomacro"
)

// Derive is a deriving state.
type Derive struct {
	macro.Context
	Macro
	Interpolator Interpolator

	first    bool
	obj      *types.TypeName
	typeSpec *ast.TypeSpec
	selector *ast.Ident
}

// NewDerive creates new Derive.
func NewDerive(m Macro) *Derive {
	selector := ast.NewIdent("m")
	d := &Derive{
		Macro:    m,
		selector: selector,
	}

	d.Interpolator = NewInterpolator(d, map[string]string{
		"$m": selector.Name,
	})
	return d
}

// With uses given context.
func (d *Derive) With(ctx macro.Context) {
	d.first = true
	d.Context = ctx
}

//nolint: unparam
func (d *Derive) impl(field Field, typ types.Type, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	call, err := d.Protocol().Impl(d, field)
	if err != nil {
		return s, err
	}

	return s.AddStmts(call), nil
}

// IsDelayed checks that given typename is marked to derive code using current macro.
// Useful for resolving cycles.
func (d *Derive) IsDelayed(typ *types.TypeName) bool {
	return d.Delayed.Find(d.Name(), typ) && !d.first
}

// IsCurrent checks that given type is current deriving type.
// Useful for resolving cycles.
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
			body, err = d.dispatch(field, typ, body)
			return body
		})

		return s, err
	}

	return d.dispatch(field, typ, s)
}

//nolint:gocyclo
func (d *Derive) dispatch(field Field, typ types.Type, s builders.StatementBuilder) (builders.StatementBuilder, error) {
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
	target := d.Macro.Target()
	if target != nil && types.Implements(typ, target) {
		return d.impl(field, target, s)
	}

	switch v := typ.(type) {
	case *types.Basic:
		if !field.Named {
			field.TypeName = nil
		}
		return d.deriveBasic(field, v, s)
	case *types.Array:
		return d.deriveArray(field, v, s)
	case *types.Slice:
		return d.deriveSlice(field, v, s)
	case *types.Struct:
		return d.deriveStruct(field, v, s)
	case *types.Pointer:
		return d.derivePtr(field, v, s)
	case *types.Interface:
		return d.deriveInterface(field, v, s)
	case *types.Map:
		if field.TypeName != nil && checkType(v.Key(), field.TypeName) {
			return s, fmt.Errorf("%w: '%s' refers to itself", ErrCycleDetected, field.TypeName.Name())
		}

		return d.deriveMap(field, v, s)
	case *types.Chan:
		return d.deriveChan(field, v, s)
	case *types.Named:
		field.TypeName = v.Obj()
		if d.IsDelayed(field.TypeName) {
			return d.impl(field, d.Macro.Target(), s)
		}
		if d.first {
			d.first = false
		}

		field.Named = true
		return d.dispatch(field, v.Underlying(), s)
	}

	return s, nil
}

// Dispatch derives code for given field.
func (d *Derive) Dispatch(field Field, typ types.Type, s builders.StatementBuilder) (builders.StatementBuilder, error) {
	return d.dispatch(field, typ, s)
}

// Derive derives code for given type declaration.
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

func (d *Derive) deriveBasic(
	field Field,
	typ *types.Basic,
	s builders.StatementBuilder,
) (builders.StatementBuilder, error) {
	if typ.Kind() == types.Invalid {
		return s, fmt.Errorf("%v: %w", field.Selector, ErrInvalidType)
	}

	block, err := d.Protocol().CallFor(d, field, typ.Kind())
	if err != nil {
		return s, err
	}

	return s.AddStmts(block), nil
}

func (d *Derive) deriveStruct(
	field Field,
	typ *types.Struct,
	s builders.StatementBuilder,
) (builders.StatementBuilder, error) {
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

func (d *Derive) tryDeriveArray(
	array Array,
	field Field,
	s builders.StatementBuilder,
) (builders.StatementBuilder, error) {
	if v, ok := d.Macro.Protocol().(ArrayDerive); ok {
		stmts, err := v.Array(d, field, array)
		if err != nil {
			return s, err
		}
		return s.AddStmts(stmts), nil
	}

	return s, nil
}

func (d *Derive) deriveArray(
	field Field,
	typ *types.Array,
	s builders.StatementBuilder,
) (builders.StatementBuilder, error) {
	return d.tryDeriveArray(Array{typ.Len(), typ.Elem()}, field, s)
}

func (d *Derive) deriveSlice(
	field Field,
	typ *types.Slice,
	s builders.StatementBuilder,
) (builders.StatementBuilder, error) {
	return d.tryDeriveArray(Array{-1, typ.Elem()}, field, s)
}

func (d *Derive) deriveInterface(
	_ Field,
	_ *types.Interface,
	s builders.StatementBuilder,
) (builders.StatementBuilder, error) {
	// Ignore interface marshaling (types which implements marshaling interface already handled)
	return s, nil
}

func (d *Derive) deriveMap(
	field Field,
	typ *types.Map,
	s builders.StatementBuilder,
) (builders.StatementBuilder, error) {
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

func (d *Derive) derivePtr(
	field Field,
	typ *types.Pointer,
	s builders.StatementBuilder,
) (builders.StatementBuilder, error) {
	if v, ok := d.Macro.Protocol().(PointerDerive); ok {
		stmts, err := v.Pointer(d, field, Pointer{
			Elem: typ.Elem(),
		})
		if err != nil {
			return s, err
		}
		return s.AddStmts(stmts), nil
	}

	return s, nil
}

func (d *Derive) deriveChan(
	field Field,
	typ *types.Chan,
	s builders.StatementBuilder,
) (builders.StatementBuilder, error) {
	return s, nil
}

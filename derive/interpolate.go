package derive

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
)

type Interpolator struct {
	derive   *Derive
	replacer *strings.Replacer
}

func NewInterpolator(derive *Derive, a ...string) Interpolator {
	return Interpolator{derive: derive, replacer: strings.NewReplacer(a...)}
}

// Interpolate interpolates given string.
func (i Interpolator) Interpolate(s string) string {
	return i.replacer.Replace(s)
}

// Expr interpolates given string and tries to parse it.
func (i Interpolator) Expr(s string) (ast.Expr, error) {
	return parser.ParseExpr(i.Interpolate(s))
}

var ErrExpected = errors.New("expected expression")

// ExprExpectKind interpolates given string, tries to parse it and checks basic type kind.
func (i Interpolator) ExprExpectKind(s string, kind types.BasicKind) (ast.Expr, error) {
	expr, err := i.Expr(s)
	if err != nil {
		return nil, err
	}

	typ, err := i.typeCheck(expr)
	if err != nil {
		return nil, err
	}

	b, ok := typ.(*types.Basic)
	if !ok || b.Kind()&kind == 0 {
		return nil, fmt.Errorf("%w type %v, got type %v '%v'", ErrExpected, kind, typ, s)
	}

	return expr, nil
}

// ExprExpectKind interpolates given string, tries to parse it and checks basic type into.
func (i Interpolator) ExprExpectInfo(s string, info types.BasicInfo) (ast.Expr, error) {
	expr, err := i.Expr(s)
	if err != nil {
		return nil, err
	}

	typ, err := i.typeCheck(expr)
	if err != nil {
		return nil, err
	}

	b, ok := typ.(*types.Basic)
	if !ok || b.Info()&info == 0 {
		return nil, fmt.Errorf("%w type %v, got type %v '%v'", ErrExpected, info, typ, s)
	}

	return expr, nil
}

func (i Interpolator) typeCheck(expr ast.Expr) (types.Type, error) {
	scope := i.derive.Package.Scope()
	scope.Insert(types.NewVar(0, i.derive.Package, "m", i.derive.obj.Type()))
	err := types.CheckExpr(token.NewFileSet(), i.derive.Package, 0, expr, i.derive.TypesInfo)
	if err != nil {
		return nil, err
	}

	typ := i.derive.TypesInfo.Types[expr]
	delete(i.derive.TypesInfo.Types, expr)

	return typ.Type, nil
}

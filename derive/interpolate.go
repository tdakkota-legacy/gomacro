package derive

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"unsafe"
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

// Copy of *types.Package
//nolint:structcheck
type typesPkg struct {
	path     string
	name     string
	scope    *types.Scope
	complete bool
	imports  []*types.Package
	fake     bool // scope lookup errors are silently dropped if package is fake (internal use only)
	cgo      bool // uses of this package will be rewritten into uses of declarations from _cgo_gotypes.go
}

func patchPackage(oldPkg *types.Package, typ types.Type) *types.Package {
	// create a copy
	pkg := types.NewPackage(oldPkg.Path(), oldPkg.Name())
	// create a child scope
	scope := types.NewScope(oldPkg.Scope(), 0, 0, "")
	// add receiver record
	scope.Insert(types.NewVar(0, oldPkg, "m", typ))
	// change scope of copy
	// TODO: find a better way to create a immutable package
	(*typesPkg)(unsafe.Pointer(pkg)).scope = scope
	return pkg
}

func (i Interpolator) patchPackage() *types.Package {
	return patchPackage(i.derive.Package, i.derive.obj.Type())
}

func (i Interpolator) typeCheck(expr ast.Expr) (types.Type, error) {
	pkg := i.patchPackage()

	err := types.CheckExpr(token.NewFileSet(), pkg, 0, expr, i.derive.TypesInfo)
	if err != nil {
		return nil, err
	}

	typ := i.derive.TypesInfo.Types[expr]
	delete(i.derive.TypesInfo.Types, expr)

	return typ.Type, nil
}

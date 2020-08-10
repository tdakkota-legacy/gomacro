package macro

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"golang.org/x/tools/go/packages"
)

func load(path ...string) ([]*packages.Package, error) {
	return packages.Load(&packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedTypesSizes | packages.NeedSyntax,
		Env:  os.Environ(),
		Fset: token.NewFileSet(),
		ParseFile: func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
			const mode = parser.AllErrors | parser.ParseComments
			return parser.ParseFile(fset, filename, src, mode)
		},
	}, path...)
}

var ErrExpectedOnlyOnePackage = errors.New("expected only one package")

func loadOne(path string) (Context, error) {
	pkgs, err := load(path)
	if err != nil {
		return Context{}, err
	}

	if len(pkgs) != 1 {
		return Context{}, ErrExpectedOnlyOnePackage
	}
	pkg := pkgs[0]

	d := Delayed{}
	d.Add(pkg)
	ctx := createContext(d, pkg)
	ctx.File = pkg.Syntax[0]

	return ctx, nil
}

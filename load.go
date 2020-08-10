package macro

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"github.com/tdakkota/gomacro/macroctx"

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

func loadOne(path string) (macroctx.Context, error) {
	pkgs, err := load(path)
	if err != nil {
		return macroctx.Context{}, err
	}

	if len(pkgs) != 1 {
		return macroctx.Context{}, ErrExpectedOnlyOnePackage
	}
	pkg := pkgs[0]

	d := macroctx.Delayed{}
	d.Add(pkg)
	ctx := createContext(d, pkg)
	ctx.File = pkg.Syntax[0]

	return ctx, nil
}

func loadDelayed(pkgs []*packages.Package) macroctx.Delayed {
	delayed := macroctx.Delayed{}
	for _, pkg := range pkgs {
		delayed.Add(pkg)
	}

	return delayed
}

func loadComments(decl ast.Decl, imports **ast.GenDecl) (comments *ast.CommentGroup) {
	switch v := decl.(type) {
	case *ast.GenDecl:
		if v.Tok == token.IMPORT {
			*imports = v
		}
		comments = v.Doc
	case *ast.FuncDecl:
		comments = v.Doc
	}

	return
}

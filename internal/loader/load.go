package loader

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/tools/go/packages"

	macro "github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/internal"
)

// LoadWalk loads packages and walks every package file.
func LoadWalk(path string, cb func(l Loaded, ctx macro.Context) error) error {
	pkgs, err := Load(path)
	if err != nil {
		return err
	}
	info := loadInfo(pkgs...)

	for _, pkg := range pkgs {
		ctx := internal.CreateContext(info.delayed, pkg)

		for _, file := range pkg.Syntax {
			ctx.File = file

			l := Loaded{Packages: info.pkgs, Module: pkg.Module}
			err := cb(l, ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Load loads packages from their path.
func Load(path string) ([]*packages.Package, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	pattern := "./..."
	loadPath := path
	if !fi.IsDir() {
		pattern = filepath.Base(path)
		loadPath = filepath.Dir(path)
	}

	return LoadPackage(loadPath, pattern, os.Environ())
}

// LoadPackage loads packages from pattern using given workdir and environment.
func LoadPackage(dir, pattern string, environ []string) ([]*packages.Package, error) {
	return packages.Load(&packages.Config{
		Dir: dir,
		Mode: packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes |
			packages.NeedSyntax |
			packages.NeedDeps,
		Env:  environ,
		Fset: token.NewFileSet(),
		ParseFile: func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
			const mode = parser.AllErrors | parser.ParseComments
			return parser.ParseFile(fset, filename, src, mode)
		},
	}, pattern)
}

// ErrExpectedOnlyOnePackage reports that loader found several packages, but expected only one.
var ErrExpectedOnlyOnePackage = errors.New("expected only one package")

// LoadOne loads exactly one package or returns error if any.
func LoadOne(path string) (macro.Context, error) {
	pkgs, err := Load(path)
	if err != nil {
		return macro.Context{}, err
	}

	if len(pkgs) != 1 {
		return macro.Context{}, ErrExpectedOnlyOnePackage
	}
	pkg := pkgs[0]

	d := macro.Delayed{}
	d.Add(pkg)
	ctx := internal.CreateContext(d, pkg)
	ctx.File = pkg.Syntax[0]

	return ctx, nil
}

// LoadReader loads source from io.Reader.
func LoadReader(r io.Reader, name string) (macro.Context, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, name, r, parser.ParseComments)
	if err != nil {
		return macro.Context{}, err
	}

	pkgName := file.Name.Name
	tpkg := types.NewPackage(pkgName, pkgName)
	info := newInfo()
	// TODO(tdakkota): it seems incorrect, to ignore errors,
	//					but there are no another way to get types
	_ = types.NewChecker(nil, fset, tpkg, info).Files([]*ast.File{file})

	arch := os.Getenv("GOARCH")
	if arch == "" {
		arch = runtime.GOARCH
	}
	pkg := &packages.Package{
		Fset:       fset,
		Types:      tpkg,
		Syntax:     []*ast.File{file},
		TypesInfo:  info,
		TypesSizes: types.SizesFor("gc", arch),
	}

	d := macro.Delayed{}
	d.Add(pkg)
	ctx := internal.CreateContext(d, pkg)
	ctx.File = file

	return ctx, nil
}

// newInfo returns a types.Info with all maps populated.
func newInfo() *types.Info {
	return &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}
}

// LoadComments gets comments from given declaration.
func LoadComments(decl ast.Decl) (comments *ast.CommentGroup) {
	switch v := decl.(type) {
	case *ast.GenDecl:
		comments = v.Doc
	case *ast.FuncDecl:
		comments = v.Doc
	}

	return
}

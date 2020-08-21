package loader

import (
	"path/filepath"

	macro "github.com/tdakkota/gomacro"
	"golang.org/x/tools/go/packages"
)

type LoadedPackages map[string]string

func (l LoadedPackages) Add(pkg *packages.Package) {
	if len(pkg.Syntax) > 0 {
		file := pkg.Fset.File(pkg.Syntax[0].Pos()).Name()
		l[pkg.Types.Path()] = filepath.Dir(file)
	}
}

func (l LoadedPackages) Has(path string) bool {
	_, ok := l[path]
	return ok
}

type info struct {
	pkgs    LoadedPackages
	delayed macro.Delayed
}

func loadInfo(pkgs ...*packages.Package) info {
	r := info{
		pkgs:    LoadedPackages{},
		delayed: macro.Delayed{},
	}

	for _, pkg := range pkgs {
		r.delayed.Add(pkg)
		r.pkgs.Add(pkg)
	}

	return r
}

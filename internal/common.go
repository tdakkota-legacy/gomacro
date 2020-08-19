package internal

import (
	"github.com/tdakkota/gomacro"
	"golang.org/x/tools/go/packages"
)

func CreateContext(delayed macro.Delayed, pkg *packages.Package) macro.Context {
	return macro.Context{
		ASTInfo: macro.ASTInfo{
			FileSet: pkg.Fset,
		},
		TypeInfo: macro.TypeInfo{
			Package:    pkg.Types,
			TypesInfo:  pkg.TypesInfo,
			TypesSizes: pkg.TypesSizes,
		},

		Delayed: delayed,
	}
}

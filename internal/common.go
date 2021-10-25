package internal

import (
	"golang.org/x/tools/go/packages"

	macro "github.com/tdakkota/gomacro"
)

// CreateContext creates new macro.Context.
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

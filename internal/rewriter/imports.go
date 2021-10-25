package rewriter

import (
	"strconv"

	"golang.org/x/tools/go/ast/astutil"

	macro "github.com/tdakkota/gomacro"
)

func (r ReWriter) fixImports(deleteUnused bool, context macro.Context) error {
	specs := astutil.Imports(context.FileSet, context.File)

	for _, group := range specs {
		for _, imprt := range group {
			importPath, _ := strconv.Unquote(imprt.Path.Value)
			if deleteUnused && !astutil.UsesImport(context.File, importPath) {
				astutil.DeleteImport(context.FileSet, context.File, importPath)
				continue
			}
		}
	}

	return nil
}

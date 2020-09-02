package rewriter

import (
	"github.com/tdakkota/gomacro/internal/rewriter/flags"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	macro "github.com/tdakkota/gomacro"
	"golang.org/x/tools/go/ast/astutil"
)

func getRelativeFilePath(base, target string) (string, error) {
	fsPath, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}

	return filepath.Rel(base, fsPath)
}

func changeImportPathBase(rel, importPath, output string) string {
	rel = strings.ReplaceAll(rel, "\\", "/")

	newPath := strings.TrimSuffix(importPath, rel) // delete subpackage path
	newPath = strings.TrimSuffix(newPath, "/")
	newPath = strings.TrimSuffix(newPath, path.Base(newPath)) // delete target path
	return path.Join(newPath, filepath.Base(output), rel)     // replace target path
}

func (r ReWriter) fixImports(deleteUnused bool, context macro.Context) error {
	specs := astutil.Imports(context.FileSet, context.File)
	absPath, err := filepath.Abs(r.source)
	if err != nil {
		return err
	}

	for _, group := range specs {
		for _, imprt := range group {
			importPath, _ := strconv.Unquote(imprt.Path.Value)
			if deleteUnused && !astutil.UsesImport(context.File, importPath) {
				astutil.DeleteImport(context.FileSet, context.File, importPath)
				continue
			}

			if !r.flags.Has(flags.AppendMode) && r.loaded.Packages.Has(importPath) {
				rel, err := getRelativeFilePath(absPath, r.loaded.Packages[importPath])
				if err != nil {
					return err
				}
				newPath := changeImportPathBase(rel, importPath, r.output)

				astutil.RewriteImport(context.FileSet, context.File, importPath, newPath)
			}
		}
	}

	return nil
}

package rewriter

import (
	"go/ast"
	"go/token"
	"os"
	"path/filepath"

	"github.com/tdakkota/gomacro"

	"golang.org/x/tools/go/packages"
)

func fixImports(imports *ast.GenDecl, file *ast.File) {
	if imports == nil {
		specs := make([]ast.Spec, len(file.Imports))
		for i := range file.Imports {
			specs[i] = file.Imports[i]
		}

		imports = &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: specs,
		}
		file.Decls = append([]ast.Decl{imports}, file.Decls...)
	} else {
		for _, imprt := range file.Imports {
			found := false
			for _, spec := range imports.Specs {
				if v, ok := spec.(*ast.ImportSpec); ok && v.Path == imprt.Path {
					found = true
					break
				}
			}

			if !found {
				imports.Specs = append(imports.Specs, imprt)
			}
		}
	}
}

func createContext(delayed macro.Delayed, pkg *packages.Package) macro.Context {
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

// prepareOutputFile creates output file directory along with any necessary parents
// and returns absolute path to output file.
// Output file path is output+filepath.Rel(path, filePath)
func prepareOutputFile(path, output, filePath string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(absPath, filePath)
	if err != nil {
		return "", err
	}

	outputFile := filepath.Join(output, rel)
	err = os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
	if err != nil {
		return "", err
	}

	return outputFile, nil
}

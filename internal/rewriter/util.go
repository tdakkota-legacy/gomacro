package rewriter

import (
	"go/ast"
	"path/filepath"
	"strings"

	builders "github.com/tdakkota/astbuilders"
)

func prepareGenFile(filePath string) (string, error) {
	ext := filepath.Ext(filePath)
	name := strings.TrimSuffix(filePath, ext)
	return name + ".gen" + ext, nil
}

func copyDecls(a []ast.Decl) []ast.Decl {
	copyDecl := make([]ast.Decl, len(a))
	copy(copyDecl, a)
	return copyDecl
}

func copyFile(file *ast.File) *ast.File {
	return builders.NewFileBuilder(file.Name.Name).
		AddImports(file.Imports...).
		Complete()
}

package macro

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/ast/astutil"
)

type ASTInfo struct {
	// Current file.
	File *ast.File
	// Current token set.
	FileSet *token.FileSet
}

// AddDecls adds declarations to current file.
func (c ASTInfo) AddDecls(decls ...ast.Decl) {
	c.File.Decls = append(c.File.Decls, decls...)
}

// AddImports adds new imports to file.
// If import already exists AddImports does nothing.
func (c ASTInfo) AddImports(newImports ...*ast.ImportSpec) {
	for _, spec := range newImports {
		path, _ := strconv.Unquote(spec.Path.Value)
		astutil.AddImport(c.FileSet, c.File, path)
	}
}

type TypeInfo struct {
	// Current package info.
	Package *types.Package
	// Type checker info.
	TypesInfo  *types.Info
	TypesSizes types.Sizes
}

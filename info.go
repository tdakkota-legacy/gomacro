package macro

import (
	"go/ast"
	"go/token"
	"go/types"
)

type ASTInfo struct {
	// Current file.
	File *ast.File
	// Current token set.
	FileSet *token.FileSet
}

type TypeInfo struct {
	// Current package info.
	Package *types.Package
	// Type checker info.
	TypesInfo  *types.Info
	TypesSizes types.Sizes
}

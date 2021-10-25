package derive

import (
	"go/ast"
	"go/types"
)

// Array is a Go slice or array type.
type Array struct {
	Size int64
	Elem types.Type
}

// ArrayDerive defines code generation for arrays and slices.
type ArrayDerive interface {
	Array(d *Derive, field Field, arr Array) (*ast.BlockStmt, error)
}

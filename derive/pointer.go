package derive

import (
	"go/ast"
	"go/types"
)

// Pointer is Go pointer type.
type Pointer struct {
	Elem types.Type
}

// PointerDerive defines code generation for pointers.
type PointerDerive interface {
	Pointer(d *Derive, field Field, p Pointer) (*ast.BlockStmt, error)
}

package derive

import (
	"go/ast"
	"go/types"
)

type Pointer struct {
	Elem types.Type
}

type PointerDerive interface {
	Pointer(d *Derive, field Field, p Pointer) (*ast.BlockStmt, error)
}

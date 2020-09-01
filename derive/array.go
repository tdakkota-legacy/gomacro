package derive

import (
	"go/ast"
	"go/types"
)

type Array struct {
	Size int64
	Elem types.Type
}

type ArrayDerive interface {
	Array(d *Derive, field Field, arr Array) (*ast.BlockStmt, error)
}

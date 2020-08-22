package derive

import (
	"go/ast"
	"go/types"
)

type Map struct {
	Key   types.Type
	Value types.Type
}

type MapDerive interface {
	Map(d *Derive, field Field, m Map) (*ast.BlockStmt, error)
}

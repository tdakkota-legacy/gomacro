package derive

import (
	"go/ast"
	"go/types"
)

// Map is a Go map type.
type Map struct {
	Key   types.Type
	Value types.Type
}

// MapDerive defines code generation for maps.
type MapDerive interface {
	Map(d *Derive, field Field, m Map) (*ast.BlockStmt, error)
}

package derive

import (
	"go/ast"
	"go/types"

	"github.com/tdakkota/gomacro/derive/base"
)

type Map struct {
	Key   types.Type
	Value types.Type
}

type MapDerive interface {
	Map(d base.Dispatcher, field base.Field, m Map) (*ast.BlockStmt, error)
}

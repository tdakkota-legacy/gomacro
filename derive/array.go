package derive

import (
	"go/ast"
	"go/types"

	"github.com/tdakkota/gomacro/derive/base"
)

type Array struct {
	Size int64
	Elem types.Type
}

type ArrayDerive interface {
	Array(d base.Dispatcher, field base.Field, arr Array) (*ast.BlockStmt, error)
}

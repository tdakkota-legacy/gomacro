package derive

import (
	"go/ast"
	"go/types"

	"github.com/tdakkota/gomacro/derive/base"
)

type Interface interface {
	CallFor(field base.Field, kind types.BasicKind) (*ast.BlockStmt, error)
	Impl(field base.Field) (*ast.BlockStmt, error)
}

type Info struct {
	Interface
	macroName string
	target    *types.Interface
}

func NewDeriveInfo(deriveInterface Interface, name string, target *types.Interface) Info {
	return Info{Interface: deriveInterface, macroName: name, target: target}
}

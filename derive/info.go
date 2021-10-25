package derive

import (
	"go/ast"
	"go/types"
)

// Protocol is an abstraction for deriving protocol.
// It defines
type Protocol interface {
	CallFor(d *Derive, field Field, kind types.BasicKind) (*ast.BlockStmt, error)
	Impl(d *Derive, field Field) (*ast.BlockStmt, error)
	Callback(d *Derive, node *ast.TypeSpec) error
}

// Macro is an abstraction for deriving macro.
type Macro interface {
	Protocol() Protocol
	Name() string
	Target() *types.Interface
}

type macroInfo struct {
	protocol Protocol
	name     string
	target   *types.Interface
}

func (m macroInfo) Protocol() Protocol {
	return m.protocol
}

func (m macroInfo) Name() string {
	return m.name
}

func (m macroInfo) Target() *types.Interface {
	return m.target
}

// CreateMacro creates derive macro from given name, target interface and Protocol.
func CreateMacro(name string, target *types.Interface, p Protocol) Macro {
	return macroInfo{
		protocol: p,
		name:     name,
		target:   target,
	}
}

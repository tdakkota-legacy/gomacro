package derive

import (
	"go/ast"

	macro "github.com/tdakkota/gomacro"
)

func Callback(m Macro) macro.HandlerFunc {
	p := m.Protocol()
	return func(cursor macro.Context, node ast.Node) error {
		if !cursor.Pre {
			if typeSpec, ok := node.(*ast.TypeSpec); ok {
				return p.Callback(NewDerive(cursor, m), typeSpec)
			}
		}

		return nil
	}
}

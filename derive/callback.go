package derive

import (
	"go/ast"

	macro "github.com/tdakkota/gomacro"
)

// Callback creates new macro handler from derive macro.
func Callback(m Macro) macro.HandlerFunc {
	p := m.Protocol()
	d := NewDerive(m)

	return func(cursor macro.Context, node ast.Node) error {
		if !cursor.Pre {
			if typeSpec, ok := node.(*ast.TypeSpec); ok {
				d.With(cursor)
				return p.Callback(d, typeSpec)
			}
		}

		return nil
	}
}

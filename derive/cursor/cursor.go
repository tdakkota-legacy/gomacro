package cursor

import (
	"go/ast"

	macro "github.com/tdakkota/gomacro"
)

func DeriveBinary(ctxt macro.Context, node ast.Node) error {
	if !ctxt.Pre {
		err := NewSerialize().Callback(ctxt, node)
		if err != nil {
			return err
		}

		err = NewDeserialize().Callback(ctxt, node)
		if err != nil {
			return err
		}
	}

	return nil
}

package cursor

import (
	"go/ast"

	"github.com/tdakkota/gomacro/derive"

	builders "github.com/tdakkota/astbuilders"
	macro "github.com/tdakkota/gomacro"
)

type DeriveBinary struct {
	Serialize   macro.Handler
	Deserialize macro.Handler
}

func Create() (*DeriveBinary, error) {
	pkgs, err := load(pkg)
	if err != nil {
		return nil, err
	}

	appender, err := target(pkgs, "Appender")
	if err != nil {
		return nil, err
	}

	scanner, err := target(pkgs, "Scanner")
	if err != nil {
		return nil, err
	}

	s := derive.CreateMacro("derive_binary", appender, &Serialize{})
	d := derive.CreateMacro("derive_binary", scanner, &Deserialize{})

	return &DeriveBinary{
		Serialize:   derive.Callback(s),
		Deserialize: derive.Callback(d),
	}, nil
}

func (d *DeriveBinary) Handle(ctxt macro.Context, node ast.Node) error {
	if !ctxt.Pre {
		err := d.Serialize.Handle(ctxt, node)
		if err != nil {
			return err
		}

		err = d.Deserialize.Handle(ctxt, node)
		if err != nil {
			return err
		}

		ctxt.AddImports(builders.Import("github.com/tdakkota/cursor"))
	}

	return nil
}

package cursor

import (
	"go/ast"

	builders "github.com/tdakkota/astbuilders"
	macro "github.com/tdakkota/gomacro"
)

type DeriveBinary struct {
	Serialize   *Serialize
	Deserialize *Deserialize
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

	return &DeriveBinary{
		Serialize:   NewSerialize(appender),
		Deserialize: NewDeserialize(scanner),
	}, nil
}

func (d *DeriveBinary) Handle(ctxt macro.Context, node ast.Node) error {
	if !ctxt.Pre {
		err := d.Serialize.Callback(ctxt, node)
		if err != nil {
			return err
		}

		err = d.Deserialize.Callback(ctxt, node)
		if err != nil {
			return err
		}

		ctxt.AddImports(builders.Import("github.com/tdakkota/cursor"))
	}

	return nil
}

package macro

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/require"
	builders "github.com/tdakkota/astbuilders"
)

func TestContext_AddImports(t *testing.T) {
	ctxt := Context{
		ASTInfo: ASTInfo{
			File: &ast.File{},
		},
	}

	imprt := builders.Import("github.com/tdakkota/astbuilders")
	ctxt.AddImports(imprt)

	a := require.New(t)
	a.Len(ctxt.File.Imports, 1)
	a.Equal(ctxt.File.Imports[0], imprt)

	ctxt.AddImports(imprt)
	a.Len(ctxt.File.Imports, 1)
	a.Equal(ctxt.File.Imports[0], imprt)

	imprt2 := builders.Import("fmt")
	ctxt.AddImports(imprt2)
	a.Len(ctxt.File.Imports, 2)
	a.Equal(ctxt.File.Imports[1], imprt2)
}

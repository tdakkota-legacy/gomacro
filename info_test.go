package macro

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdakkota/astbuilders"
)

func TestASTInfo_AddDecls(t *testing.T) {
	ctxt := Context{
		ASTInfo: ASTInfo{
			File: &ast.File{},
		},
	}

	f := &ast.FuncDecl{
		Name: ast.NewIdent("f"),
		Type: builders.FuncType()(),
	}
	ctxt.AddDecls(f)

	require.Equal(t, f, ctxt.File.Decls[0])
}

func TestASTInfo_AddImports(t *testing.T) {
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

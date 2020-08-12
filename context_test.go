package macro

import (
	builders "github.com/tdakkota/astbuilders"
	"go/ast"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_importEqual(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		a, b := &ast.ImportSpec{
			Name: ast.NewIdent("a"),
			Path: &ast.BasicLit{Value: `"import"`},
		}, &ast.ImportSpec{
			Name: ast.NewIdent("a"),
			Path: &ast.BasicLit{Value: `"import"`},
		}

		require.True(t, importEqual(a, b))
	})

	t.Run("different-name", func(t *testing.T) {
		a, b := &ast.ImportSpec{
			Name: ast.NewIdent("a"),
			Path: &ast.BasicLit{Value: `"import"`},
		}, &ast.ImportSpec{
			Name: ast.NewIdent("b"),
			Path: &ast.BasicLit{Value: `"import"`},
		}

		require.False(t, importEqual(a, b))
	})

	t.Run("different-path", func(t *testing.T) {
		a, b := &ast.ImportSpec{
			Name: ast.NewIdent("a"),
			Path: &ast.BasicLit{Value: `"import"`},
		}, &ast.ImportSpec{
			Name: ast.NewIdent("a"),
			Path: &ast.BasicLit{Value: `"import2"`},
		}

		require.False(t, importEqual(a, b))
	})
}

func TestContext_AddImports(t *testing.T) {
	ctxt := Context{
		File: &ast.File{},
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

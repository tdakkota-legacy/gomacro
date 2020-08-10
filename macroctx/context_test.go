package macroctx

import (
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

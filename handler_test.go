package macro

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tdakkota/gomacro/macroctx"
)

func TestOnlyFunction(t *testing.T) {
	counter := 0
	handler := OnlyFunction("callme", func(ctx macroctx.Context, call *ast.CallExpr) error {
		counter++
		return nil
	})

	err := handler.Handle(macroctx.Context{}, &ast.CallExpr{Fun: ast.NewIdent("callme")})
	require.NoError(t, err)
	require.Equal(t, 1, counter)

	err = handler.Handle(macroctx.Context{}, &ast.CallExpr{Fun: ast.NewIdent("notcallme")})
	require.NoError(t, err)
	require.Equal(t, 1, counter)
}

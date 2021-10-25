package deepcopy

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_changeSelectorHead(t *testing.T) {
	a := require.New(t)
	expr, err := parser.ParseExpr("(*a).b[0].call().c")
	a.NoError(err)

	expr, idx, err := changeSelectorHead(expr, ast.NewIdent("r"))
	a.NoError(err)
	a.Equal(1, idx)

	var s strings.Builder
	a.NoError(printer.Fprint(&s, token.NewFileSet(), expr))
	a.Equal("(*r).b[0].call().c", s.String())
}

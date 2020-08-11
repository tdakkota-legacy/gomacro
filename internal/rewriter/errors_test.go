package rewriter

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorPrinter(t *testing.T) {
	type report = struct {
		pos token.Pos
		err error
	}

	var stack []report
	cb := func(pos token.Pos, err error) {
		stack = append(stack, report{pos, err})
	}

	printer := errorPrinter{cb}
	printer.Reportf(0, "%d", 1)
	require.Equal(t, token.Pos(0), stack[0].pos)
	require.Equal(t, "1", stack[0].err.Error())

	printer.Report(1, "hi")
	require.Equal(t, token.Pos(1), stack[1].pos)
	require.Equal(t, "hi", stack[1].err.Error())
}

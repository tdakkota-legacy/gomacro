package macro

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext_Reportf(t *testing.T) {
	r := Report{}
	ctxt := Context{
		Report: func(report Report) {
			r = report
		},
	}

	ctxt.Reportf(11, "%d", 10)
	require.Equal(t, r.Pos, token.Pos(11))
	require.Equal(t, r.Message, "10")
}

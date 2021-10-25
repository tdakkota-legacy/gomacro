package runner

import (
	"testing"

	"github.com/stretchr/testify/require"

	macro "github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/internal/rewriter"
)

func TestRunner_run(t *testing.T) {
	t.Run("check-path-replace", func(t *testing.T) {
		r := Runner{}
		err := r.run(map[string]macro.Handler{}, func(writer rewriter.ReWriter) error {
			require.Equal(t, "./", writer.Source())
			return nil
		})
		require.NoError(t, err)
	})

	t.Run("err-stop", func(t *testing.T) {
		r := Runner{}
		err := r.run(map[string]macro.Handler{}, func(writer rewriter.ReWriter) error {
			return macro.ErrStop
		})
		require.NoError(t, err)
	})
}

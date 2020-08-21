package rewriter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_urlRel(t *testing.T) {
	t.Run("rel", func(t *testing.T) {
		path := "github.com/a/b/c"
		s, err := urlRel("github.com/a", path)
		require.NoError(t, err)
		require.Equal(t, "/b/c", s)
	})

	t.Run("same", func(t *testing.T) {
		path := "github.com/a/b/c"
		s, err := urlRel(path, path)
		require.NoError(t, err)
		require.Equal(t, path, s)
	})

	t.Run("not-rel", func(t *testing.T) {
		_, err := urlRel("github.com/a/b/c", "gitlab.com/a/b/c")
		require.Error(t, err)
	})

}

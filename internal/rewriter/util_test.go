package rewriter

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_urlRel(t *testing.T) {
	testPath := "github.com/a/b/c"
	t.Run("rel", func(t *testing.T) {
		s, err := urlRel("github.com/a", testPath)
		require.NoError(t, err)
		require.Equal(t, "/b/c", s)
	})

	t.Run("same", func(t *testing.T) {
		s, err := urlRel(testPath, testPath)
		require.NoError(t, err)
		require.Equal(t, testPath, s)
	})

	t.Run("not-rel", func(t *testing.T) {
		_, err := urlRel(testPath, strings.ReplaceAll(testPath, "github", "gitlab"))
		require.Error(t, err)
	})

}

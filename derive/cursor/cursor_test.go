package cursor

import (
	"go/ast"
	"os"
	"path/filepath"
	"testing"

	"github.com/tdakkota/gomacro/internal/testutil"
	"github.com/tdakkota/gomacro/runner"

	macro "github.com/tdakkota/gomacro"
)

func DeriveBinary(ctxt macro.Context, node ast.Node) error {
	if !ctxt.Pre {
		err := NewSerialize().Callback(ctxt, node)
		if err != nil {
			return err
		}

		err = NewDeserialize().Callback(ctxt, node)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestCursor(t *testing.T) {
	err := testutil.WithTempDir("cursor-test", func(path string) error {
		outputFile := filepath.Join(path, "cursor.go")

		err := runner.Run("./testdata/src/cursor.go", outputFile, macro.Macros{
			"derive_binary": macro.HandlerFunc(DeriveBinary),
		})
		if err != nil {
			return err
		}

		err = testutil.RunGoTool(os.Stdout, "run", outputFile)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}

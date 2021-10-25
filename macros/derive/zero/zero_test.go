package zero

import (
	"os"
	"path/filepath"
	"testing"

	macro "github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/internal/testutil"
	"github.com/tdakkota/gomacro/runner"
)

func TestZero(t *testing.T) {
	err := testutil.WithTempDir("zero-test", func(path string) error {
		outputFile := filepath.Join(path, "zero.go")

		r := runner.Runner{
			Source: "./testdata/zero.go",
			Output: outputFile,
		}
		err := r.Run(macro.Macros{
			"derive_zero": Zero("derive_zero"),
		})
		if err != nil {
			return err
		}

		if err := testutil.CopyFile(filepath.Join(path, "gen.go"), r.Source); err != nil {
			return err
		}

		err = testutil.RunGoTool(os.Stdout, "run", outputFile, filepath.Join(path, "gen.go"))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}

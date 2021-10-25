package deepcopy

import (
	"os"
	"path/filepath"
	"testing"

	macro "github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/internal/testutil"
	"github.com/tdakkota/gomacro/runner"
)

func TestCopy(t *testing.T) {
	err := testutil.WithTempDir("copy-test", func(path string) error {
		outputFile := filepath.Join(path, "copy.go")

		r := runner.Runner{
			Source: "./testdata/copy.go",
			Output: outputFile,
		}
		err := r.Run(macro.Macros{
			"derive_copy": DeepCopy("derive_copy"),
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
